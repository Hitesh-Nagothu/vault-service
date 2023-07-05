package service

import (
	"errors"
	"io"
	"mime"
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

	fileTypes := fs.GetFileType(fileHeader)

	if len(fileTypes) == 0 {
		fs.logger.Error("Failed to infer the type of file uploaded")
		return errors.New("failed to infer the type of file uploaded")
	}

	fileType, isAllowed := fs.IsAllowedFileType(fileTypes)
	if !isAllowed {
		fs.logger.Error("Invalid file type", zap.String("requested_file_type", fileType), zap.Strings("allowed_types", fileTypes))
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
	_, createFileErr := fs.repo.Add(newFile)
	if createFileErr != nil {
		fs.logger.Error("Failed to create new file. Aborting file upload")
		return errors.New("something went wrong processing the file")
	}

	//get ther user uploading the file
	user, err := fs.userService.GetUser(userEmail)
	if err != nil || utility.IsStructEmpty(user) {
		fs.logger.Error("Failed to find the user owner for file. Aborting file upload")
		return errors.New("something went wrong processing the file")
	}

	//TODO
	//update the existing user with new file

	return nil
}

func (fs *FileService) GetFileType(fileHeader *multipart.FileHeader) []string {
	contentType := fileHeader.Header.Get("Content-Type")
	extension, _ := mime.ExtensionsByType(contentType)
	extensionsWithoutPrefix := make([]string, len(extension))

	if len(extension) > 0 {
		for _, ext := range extension {
			extensionsWithoutPrefix = append(extensionsWithoutPrefix, strings.TrimPrefix(ext, "."))
		}
		return extensionsWithoutPrefix
	}
	return make([]string, 0)
}

func (fs *FileService) IsAllowedFileType(fileTypes []string) (string, bool) {

	for _, fileType := range fileTypes {
		switch fileType {
		case "jpg", "jpeg", "png", "gif", "pdf", "txt", "doc", "docx":
			return fileType, true
		default:
			return "", false
		}
	}

	return "", false
}

func (fs *FileService) GetIPFSHashForFile(fileData []byte) (string, error) {

	ipfsInstance := fs.ipfsService.GetIPFSInstance()
	hash, err := ipfsInstance.AddContent(fileData)
	if err != nil {
		//TODO
	}

	return hash, nil
}
