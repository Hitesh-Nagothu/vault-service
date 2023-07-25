package service

import (
	"errors"
	"io"
	"mime/multipart"
	"strings"

	"github.com/Hitesh-Nagothu/vault-service/data"
	"github.com/Hitesh-Nagothu/vault-service/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type FileService struct {
	repo         *data.FileRepository
	logger       *zap.Logger
	ipfsService  *IPFSService
	chunkService *ChunkService
	userService  *UserService
}

func NewFileService(logger *zap.Logger, repo *data.FileRepository, ipfsService *IPFSService, chunkService *ChunkService, userService *UserService) *FileService {
	return &FileService{
		logger:       logger,
		repo:         repo,
		ipfsService:  ipfsService,
		chunkService: chunkService,
		userService:  userService,
	}
}

const (
	MaxFileSize = 5 * 1024 * 1024 // 5MB in bytes
)

func (fs *FileService) CreateFile(file multipart.File, fileHeader *multipart.FileHeader, userEmail string) error {

	if fileHeader.Size > MaxFileSize {
		return errors.New("file size uploaded exceeds the permissible limit of 5MB")
	}

	fileType := fs.GetFileType(fileHeader.Filename)
	fileType, isAllowed := fs.IsAllowedFileType(fileType)
	if !isAllowed {
		fs.logger.Error("Invalid file type", zap.String("requested_file_type", fileType))
		return errors.New("unsupported file type uploaded")
	}

	filebytes, readErr := io.ReadAll(file)
	if readErr != nil {
		fs.logger.Error("Failed to read file", zap.Error(readErr))
		return errors.New("something went wrong reading the file")
	}

	fileHash, hashGenErr := fs.GetIPFSHashForFile(filebytes)
	if hashGenErr != nil {
		fs.logger.Error("Failed to generate hash for file", zap.String("fileName", fileHeader.Filename), zap.String("fileType", fileType))
		return errors.New("something went wrong processing the file")
	}

	//Note: Following only a single chunk for a file, will likely add the chunk implementation later if needed

	//insert the new chunk
	createdChunk, createChunkErr := fs.chunkService.CreateChunk(fileHash)
	if createChunkErr != nil {
		fs.logger.Error("Failed to get a new chunk for file. Aborting file upload")
		return errors.New("something went wrong processing the file")
	}

	newFile := data.File{
		Name: fileHeader.Filename,
		Type: fileType,
		ChunkIDs: []primitive.ObjectID{
			createdChunk.ID,
		},
	}

	//insert the new file
	createdFile, createFileErr := fs.repo.Add(newFile)
	if createFileErr != nil {
		fs.logger.Error("Failed to create new file. Aborting file upload")
		return errors.New("something went wrong processing the file")
	}

	//get ther user uploading the file
	user, err := fs.userService.GetUser(userEmail)
	//TODO check for notfound error else abort
	if err != nil || utility.IsStructEmpty(user) {
		user, err = fs.userService.CreateUser(userEmail)
		fs.logger.Info("Create a new user previously not found", zap.String("user_email", user.Email))
	}

	userUpdate := data.User{
		Files: []primitive.ObjectID{createdFile.ID}, //sending partial object
	}

	updateUserErr := fs.userService.UpdateUser(user.ID, userUpdate)
	if updateUserErr != nil {
		fs.logger.Error("Failed to udpate the user owner with new file. Aborting file upload")
		return errors.New("something went wrong processing the file")
	}

	//TODO make chunk storing, file creation and update user with new file transactional
	fs.logger.Info("File upload successful", zap.String("file_name", createdFile.Name))
	return nil
}

func (fs *FileService) GetFileType(filename string) string {
	return strings.Split(filename, ".")[1]
}

func (fs *FileService) IsAllowedFileType(fileType string) (string, bool) {

	switch fileType {
	case "jpg", "jpeg", "png", "gif", "pdf", "txt", "doc", "docx":
		return fileType, true
	default:
		return "", false
	}

}

func (fs *FileService) GetIPFSHashForFile(fileData []byte) (string, error) {

	ipfsInstance := fs.ipfsService.GetIPFSInstance()
	hash, err := ipfsInstance.AddContent(fileData)
	if err != nil {
		//TODO
	}
	return hash, nil
}
