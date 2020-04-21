package proj2

// CS 161 Project 2 Spring 2020
// You MUST NOT change what you import.  If you add ANY additional
// imports it will break the autograder. We will be very upset.

import (
	// You neet to add with
	// go get github.com/cs161-staff/userlib
	"github.com/cs161-staff/userlib"

	// Life is much easier with json:  You are
	// going to want to use this so you can easily
	// turn complex structures into strings etc...
	"encoding/json"

	// Likewise useful for debugging, etc...
	"encoding/hex"

	// UUIDs are generated right based on the cryptographic PRNG
	// so lets make life easier and use those too...
	//
	// You need to add with "go get github.com/google/uuid"
	"github.com/google/uuid"

	// Useful for debug messages, or string manipulation for datastore keys.
	"strings"

	// Want to import errors.
	"errors"

	// Optional. You can remove the "_" there, but please do not touch
	// anything else within the import bracket.
	_ "strconv"

	// if you are looking for fmt, we don't give you fmt, but you can use userlib.DebugMsg.
	// see someUsefulThings() below:
)

// This serves two purposes:
// a) It shows you some useful primitives, and
// b) it suppresses warnings for items not being imported.
// Of course, this function can be deleted.
func someUsefulThings() {
	// Creates a random UUID
	userlib.SetDebugStatus(true)
	f := uuid.New()
	userlib.DebugMsg("UUID as string:%v", f.String())

	// Example of writing over a byte of f
	f[0] = 10
	userlib.DebugMsg("UUID as string:%v", f.String())

	// takes a sequence of bytes and renders as hex
	h := hex.EncodeToString([]byte("fubar"))
	userlib.DebugMsg("The hex: %v", h)

	// Marshals data into a JSON representation
	// Will actually work with go structures as well
	d, _ := json.Marshal(f)
	userlib.DebugMsg("The json data: %v", string(d))
	var g uuid.UUID
	json.Unmarshal(d, &g)
	userlib.DebugMsg("Unmashaled data %v", g.String())

	// This creates an error type
	userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("This is an error")))

	// And a random RSA key.  In this case, ignoring the error
	// return value
	var pk userlib.PKEEncKey
        var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("Key is %v, %v", pk, sk)
}

// Helper function: Takes the first 16 bytes and
// converts it into the UUID type
func bytesToUUID(data []byte) (ret uuid.UUID) {
	for x := range ret {
		ret[x] = data[x]
	}
	return
}

func createEncryptionKeys() (ret EncryptionData) {
	var key EncryptionData
	key.RecordLocator = uuid.New()
	key.SymmetricKey = userlib.RandomBytes(16)
	key.MACKey = userlib.RandomBytes(16)
	return key
}

func createUnlockInfo() (ret UnlockInfo) {
	var key UnlockInfo
	key.DSKey = uuid.New()
	key.DecryptionKey = userlib.RandomBytes(16)
	key.MACKey = userlib.RandomBytes(16)
	return key
}

func createSharedMDInfo() (ret SharedMetaDataInfo) {
	var key SharedMetaDataInfo
	key.DSKey = uuid.New()
	key.DecryptionKey = userlib.RandomBytes(16)
	key.MACKey = userlib.RandomBytes(16)
	return key
}

func encryptAndStore(data []byte, enc EncryptionData) (e error) {
	encData := userlib.SymEnc(enc.SymmetricKey, userlib.RandomBytes(16), data)
	macTag, err := userlib.HMACEval(enc.MACKey, encData)
	if err != nil{
		return err
	}
	finalData := append(macTag[:], encData[:]...)
	userlib.DatastoreSet(enc.RecordLocator, finalData)
	return nil
}

func fetchAndDecrypt(dec EncryptionData) (ret []byte, e error) {
	var encData []byte
	var ok bool
	if encData, ok = userlib.DatastoreGet(dec.RecordLocator); !ok {
		return nil, errors.New("Invalid credentials")
	}
	if len(encData) < 64 {
		return nil, errors.New("File has been tampered with")
	}
	checkMac, err := userlib.HMACEval(dec.MACKey, encData[64:])
	if err != nil {
		return nil, err
	}
	if !userlib.HMACEqual(checkMac, encData[:64]) {
		return nil, errors.New("File has been tampered with")
	}
	decData := userlib.SymDec(dec.SymmetricKey, encData[64:])
	return decData, nil
}

func storeUser(userData *User) (e error) {
	byteData, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	return encryptAndStore(byteData, userData.UserEnc)
}


func ceil(lnt int) (x int) {
	int_div := lnt / 1024
	dec_div := float64(lnt) / 1024.0
	if float64(int_div) != dec_div {
		return int_div + 1
	}
	return int_div
}

func getMetaData(info SharedMetaDataInfo, deletePrevious bool) (md MetaData, e error) {
	var blankMd MetaData
	encData, ok := userlib.DatastoreGet(info.DSKey)
	if !ok {
		return blankMd, errors.New("File is nonexistant")
	}
	if len(encData) < 64 {
		return blankMd, errors.New("File has been tampered with")
	}
	checkMac, _ := userlib.HMACEval(info.MACKey, encData[64:])
	if !userlib.HMACEqual(encData[:64], checkMac) {
		return blankMd, errors.New("File has been tampered with")
	}
	var ret MetaData
	_ = json.Unmarshal(userlib.SymDec(info.DecryptionKey, encData[64:]), &ret)
	if deletePrevious {
		userlib.DatastoreDelete(info.DSKey)
	}
	return ret, nil
}

func storeMetaData(FileData MetaData, info SharedMetaDataInfo) (e error) {
	packaged_data, _ := json.Marshal(FileData)
	encInfo := userlib.SymEnc(info.DecryptionKey, userlib.RandomBytes(16), packaged_data)
	macTag, _ := userlib.HMACEval(info.MACKey, encInfo)
	encInfo = append(macTag[:], encInfo[:]...)
	userlib.DatastoreSet(info.DSKey, encInfo)
	return nil
}

func storeMetaDataInfo(metaDataInfo SharedMetaDataInfo, info UnlockInfo) (e error) {
	packaged_data, _ := json.Marshal(metaDataInfo)
	encInfo := userlib.SymEnc(info.DecryptionKey, userlib.RandomBytes(16), packaged_data)
	macTag, _ := userlib.HMACEval(info.MACKey, encInfo)
	encInfo = append(macTag[:], encInfo[:]...)
	userlib.DatastoreSet(info.DSKey, encInfo)
	return nil
}

func getMetaDataInfo(info UnlockInfo) (mdi SharedMetaDataInfo, e error) {
	var blankMdi SharedMetaDataInfo
	encData, ok := userlib.DatastoreGet(info.DSKey)
	if !ok {
		return blankMdi, errors.New("File nonexistant")
	}
	if len(encData) < 64 {
		return blankMdi, errors.New("File has been tampered with")
	}
	checkMac, _ := userlib.HMACEval(info.MACKey, encData[64:])
	if !userlib.HMACEqual(encData[:64], checkMac) {
		return blankMdi, errors.New("File has been tampered with")
	}
	var ret SharedMetaDataInfo
	_ = json.Unmarshal(userlib.SymDec(info.DecryptionKey, encData[64:]), &ret)
	return ret, nil
}

// The structure definition for a user record
type User struct {
	Username string
	UserEnc EncryptionData
	FileListEnc EncryptionData
	SharedFileListEnc EncryptionData
	RSAPrivKey userlib.PKEDecKey
	SigPrivKey userlib.DSSignKey
}

type EncryptionData struct {
	RecordLocator uuid.UUID
	SymmetricKey []byte
	MACKey []byte
}

type MetaData struct {
	NumBlocks int
	FileKey []byte
	FileMACKey []byte
	LastBlock []byte
	LastMACTag []byte
}

type OwnedFileMetaDataInfo struct {
	OriginalInfo SharedMetaDataInfo
	SharedUnlocks map[string]UnlockInfo
}

type SharedMetaDataInfo struct {
	DSKey uuid.UUID
	MACKey []byte
	DecryptionKey []byte
}

type UnlockInfo struct {
	DSKey uuid.UUID
	MACKey []byte
	DecryptionKey []byte
}

	// You can add other fields here if you want...
	// Note for JSON to marshal/unmarshal, the fields need to
	// be public (start with a capital letter)


// This creates a user.  It will only be called once for a user
// (unless the keystore and datastore are cleared during testing purposes)

// It should store a copy of the userdata, suitably encrypted, in the
// datastore and should store the user's public key in the keystore.

// The datastore may corrupt or completely erase the stored
// information, but nobody outside should be able to get at the stored

// You are not allowed to use any global storage other than the
// keystore and the datastore functions in the userlib library.

// You can assume the password has strong entropy, EXCEPT
// the attackers may possess a precomputed tables containing
// hashes of common passwords downloaded from the internet.
func InitUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	if _, ok := userlib.KeystoreGet(username + "/PK"); ok {
		return &userdata, errors.New("User already exists")
	}
	userEncryptionInfo := userlib.Argon2Key([]byte(password), []byte(username), 48)
	userdataptr = &userdata
	userdata.Username = username
	var userKey EncryptionData
	userKey.RecordLocator, _ = uuid.FromBytes(userEncryptionInfo[:16])
	userKey.SymmetricKey = userEncryptionInfo[16:32]
	userKey.MACKey = userEncryptionInfo[32:]
	userdata.UserEnc = userKey
	userdata.FileListEnc = createEncryptionKeys()
	userdata.SharedFileListEnc = createEncryptionKeys()
	fileList := make(map[string]OwnedFileMetaDataInfo)
	sharedFileList := make(map[string]UnlockInfo)
	fDt, _ := json.Marshal(fileList)
	encryptAndStore(fDt, userdata.FileListEnc)
	sDt, _ := json.Marshal(sharedFileList)
	encryptAndStore(sDt, userdata.SharedFileListEnc)
	var publicKey userlib.PKEEncKey
	publicKey, userdata.RSAPrivKey, _ = userlib.PKEKeyGen()
	var signPublicKey userlib.DSVerifyKey
	userdata.SigPrivKey, signPublicKey, _ = userlib.DSKeyGen()
	_ = userlib.KeystoreSet(username + "/PK", publicKey)
	_ = userlib.KeystoreSet(username + "/DS", signPublicKey)
	err = storeUser(&userdata)
	if err != nil{
		return &userdata, err
	}
	// someUsefulThings()
	return &userdata, nil
}

// This fetches the user information from the Datastore.  It should
// fail with an error if the user/password is invalid, or if the user
// data was corrupted, or if the user can't be found.
func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	userEncryptionInfo := userlib.Argon2Key([]byte(password), []byte(username), 48)
	var userKey EncryptionData
	userKey.RecordLocator, _ = uuid.FromBytes(userEncryptionInfo[:16])
	userKey.SymmetricKey = userEncryptionInfo[16:32]
	userKey.MACKey = userEncryptionInfo[32:]
	userData, err := fetchAndDecrypt(userKey)
	if err != nil {
		return &userdata, err
	}
	err = json.Unmarshal(userData, &userdata)
	if err != nil {
		return &userdata, err
	}
	if userdata.Username != username {
		return nil, errors.New("File has been tampered with")
	}
	return &userdata, nil
}

// This stores a file in the datastore.
//
// The plaintext of the filename + the plaintext and length of the filename
// should NOT be revealed to the datastore!

func (userdata *User) StoreFile(filename string, data []byte) {
	var FileData MetaData
	var err error
	var metaDataInfo SharedMetaDataInfo
	var fileList map[string]OwnedFileMetaDataInfo
	var sharedFileList map[string]UnlockInfo
	fD, err1 := fetchAndDecrypt(userdata.FileListEnc)
	if err1 != nil {
		return
	}
	json.Unmarshal(fD, &fileList)
	sD, err2 := fetchAndDecrypt(userdata.SharedFileListEnc)
	if err2 != nil {
		return
	}
	json.Unmarshal(sD, &sharedFileList)
	if ownedInfo, ok := fileList[filename]; ok {
		metaDataInfo = ownedInfo.OriginalInfo
		getMetaData(metaDataInfo, true)
	} else if unlockInfo, ok2 := sharedFileList[filename]; ok2 {
		metaDataInfo, err = getMetaDataInfo(unlockInfo)
		if err != nil {
			return
		}
		FileData, err = getMetaData(metaDataInfo, true)
		if err != nil {
			return
		}
	} else {
		metaDataInfo = createSharedMDInfo()
		var newOwnedInfo OwnedFileMetaDataInfo
		newOwnedInfo.OriginalInfo = metaDataInfo
		newOwnedInfo.SharedUnlocks = make(map[string]UnlockInfo)
		fileList[filename] = newOwnedInfo
	}
	FileData.NumBlocks = 1
	FileData.FileKey = userlib.RandomBytes(16)
	FileData.FileMACKey = userlib.RandomBytes(16)
	FileData.LastBlock = userlib.RandomBytes(16)
	data = append(userlib.RandomBytes(80), data[:]...)
	encBlock := userlib.SymEnc(FileData.FileKey, userlib.RandomBytes(userlib.AESBlockSize), data)
	FileData.LastMACTag, _ = userlib.HMACEval(FileData.FileMACKey, encBlock)
	k, _ := uuid.FromBytes(FileData.LastBlock)
	userlib.DatastoreSet(k, encBlock)
	storeMetaData(FileData, metaDataInfo)
	byteData, _ := json.Marshal(fileList)
	encryptAndStore(byteData, userdata.FileListEnc)
	return
}

// This adds on to an existing file.
//
// Append should be efficient, you shouldn't rewrite or reencrypt the
// existing file, but only whatever additional information and
// metadata you need.
func (userdata *User) AppendFile(filename string, data []byte) (err error) {
	var FileData MetaData
	var metaDataInfo SharedMetaDataInfo
	var fileList map[string]OwnedFileMetaDataInfo
	var sharedFileList map[string]UnlockInfo
	fD, err1 := fetchAndDecrypt(userdata.FileListEnc)
	if err1 != nil {
		return err1
	}
	json.Unmarshal(fD, &fileList)
	sD, err2 := fetchAndDecrypt(userdata.SharedFileListEnc)
	if err2 != nil {
		return err1
	}
	json.Unmarshal(sD, &sharedFileList)
	if ownedInfo, ok := fileList[filename]; ok {
		metaDataInfo = ownedInfo.OriginalInfo
		FileData, err = getMetaData(metaDataInfo, false)
		if err != nil {
			return err
		}
	} else if unlockInfo, ok2 := sharedFileList[filename]; ok2 {
		metaDataInfo, err = getMetaDataInfo(unlockInfo)
		if err != nil {
			return err
		}
		FileData, err = getMetaData(metaDataInfo, false)
		if err != nil {
			return err
		}
	} else {
		return errors.New("File is neither owned or shared, nonexistant")
	}
	data = append(FileData.LastMACTag[:], data[:]...)
	data = append(FileData.LastBlock[:], data[:]...)
	FileData.NumBlocks += 1
	FileData.LastBlock = userlib.RandomBytes(16)
	encBlock := userlib.SymEnc(FileData.FileKey, userlib.RandomBytes(userlib.AESBlockSize), data)
	FileData.LastMACTag, _ = userlib.HMACEval(FileData.FileMACKey, encBlock)
	k, _ := uuid.FromBytes(FileData.LastBlock)
	userlib.DatastoreSet(k, encBlock)
	storeMetaData(FileData, metaDataInfo)
	return nil
}

// This loads a file from the Datastore.
//
// It should give an error if the file is corrupted in any way.
func (userdata *User) LoadFile(filename string) (data []byte, err error) {
	var FileData MetaData
	var metaDataInfo SharedMetaDataInfo
	var fileList map[string]OwnedFileMetaDataInfo
	var sharedFileList map[string]UnlockInfo
	fD, err1 := fetchAndDecrypt(userdata.FileListEnc)
	if err1 != nil {
		return nil, err1
	}
	json.Unmarshal(fD, &fileList)
	sD, err2 := fetchAndDecrypt(userdata.SharedFileListEnc)
	if err2 != nil {
		return nil, err2
	}
	json.Unmarshal(sD, &sharedFileList)
	if ownedInfo, ok := fileList[filename]; ok {
		metaDataInfo = ownedInfo.OriginalInfo
		FileData, err = getMetaData(metaDataInfo, false)
		if err != nil {
			return nil, err
		}
	} else if unlockInfo, ok2 := sharedFileList[filename]; ok2 {
		metaDataInfo, err = getMetaDataInfo(unlockInfo)
		if err != nil {
			return nil, err
		}
		FileData, err = getMetaData(metaDataInfo, false)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("File is neither owned or shared, nonexistant")
	}
	var ret []byte
	curKey := FileData.LastBlock
	curTag := FileData.LastMACTag
	for i := 1; i <= FileData.NumBlocks; i += 1 {
		k, _ := uuid.FromBytes(curKey)
		block, ok := userlib.DatastoreGet(k)
		if !ok {
			return nil, errors.New("File does not exist")
		}
		checkMac, _ := userlib.HMACEval(FileData.FileMACKey, block)
		if !userlib.HMACEqual(checkMac, curTag) {
			return nil, errors.New("File has been tampered with")
		}
		decryptedBlock := userlib.SymDec(FileData.FileKey, block)
		if len(decryptedBlock) < 80 {
			return nil, errors.New("File has been tampered with")
		}
		curKey = decryptedBlock[:16]
		curTag = decryptedBlock[16:80]
		ret = append(decryptedBlock[80:], ret[:]...)
	}
	return ret, nil
}

// This creates a sharing record, which is a key pointing to something
// in the datastore to share with the recipient.

// This enables the recipient to access the encrypted file as well
// for reading/appending.

// Note that neither the recipient NOR the datastore should gain any
// information about what the sender calls the file.  Only the
// recipient can access the sharing record, and only the recipient
// should be able to know the sender.
func (userdata *User) ShareFile(filename string, recipient string) (
	magic_string string, err error) {
	var fileList map[string]OwnedFileMetaDataInfo
	var sharedFileList map[string]UnlockInfo
	fD, err1 := fetchAndDecrypt(userdata.FileListEnc)
	if err1 != nil {
		return "", err1
	}
	json.Unmarshal(fD, &fileList)
	sD, err2 := fetchAndDecrypt(userdata.SharedFileListEnc)
	if err2 != nil {
		return "", err2
	}
	json.Unmarshal(sD, &sharedFileList)
	retKey := uuid.New()
	var keyInfo UnlockInfo
	if ownedInfo, ok := fileList[filename]; ok {
		keyInfo = createUnlockInfo()
		storeMetaDataInfo(ownedInfo.OriginalInfo, keyInfo)
		ownedInfo.SharedUnlocks[recipient] = keyInfo
	} else if unlockInfo, ok2 := sharedFileList[filename]; ok2 {
		keyInfo = unlockInfo
		mdi, e1 := getMetaDataInfo(keyInfo)
		if e1 != nil {
			return "", errors.New("File has been tampered with")
		}
		_, e2 := getMetaData(mdi, false)
		if e2 != nil {
			return "", errors.New("No permission to share this file")
		}
	} else {
		return "", errors.New("No such file")
	}
	byteUUID, _ := json.Marshal(keyInfo.DSKey)
	packaged_data := append(byteUUID[:], keyInfo.MACKey[:]...)
	packaged_data = append(packaged_data[:], keyInfo.DecryptionKey[:]...)
	recipientKey, okKey := userlib.KeystoreGet(recipient + "/PK")
	if !okKey {
		return "", errors.New("No such recipient")
	}
	encryptedKeyInfo, pkeerr := userlib.PKEEnc(recipientKey, packaged_data)
	if pkeerr != nil {
		return "", pkeerr
	}
	signature, _ := userlib.DSSign(userdata.SigPrivKey, encryptedKeyInfo)
	encryptedKeyInfo = append(signature[:], encryptedKeyInfo[:]...)
	userlib.DatastoreSet(retKey, encryptedKeyInfo)
	ret, _ := json.Marshal(retKey)
	byteData, _ := json.Marshal(fileList)
	encryptAndStore(byteData, userdata.FileListEnc)
	return string(ret), nil
}

// Note recipient's filename can be different from the sender's filename.
// The recipient should not be able to discover the sender's view on
// what the filename even is!  However, the recipient must ensure that
// it is authentically from the sender.
func (userdata *User) ReceiveFile(filename string, sender string,
	magic_string string) error {
	var sharedFileList map[string]UnlockInfo
	sD, err2 := fetchAndDecrypt(userdata.SharedFileListEnc)
	if err2 != nil {
		return err2
	}
	json.Unmarshal(sD, &sharedFileList)
	if _, ok := sharedFileList[filename]; ok {
		return errors.New("Already have a file under this name")
	}
	var dsKey uuid.UUID
	_ = json.Unmarshal([]byte(magic_string), &dsKey)
	info, ok := userlib.DatastoreGet(dsKey)
	if !ok {
		return errors.New("No such file")
	}
	if len(info) < 256 {
		return errors.New("File has been tampered with")
	}
	verifyKey, ok2 := userlib.KeystoreGet(sender + "/DS")
	if !ok2 {
		return errors.New("No such sender")
	}
	err := userlib.DSVerify(verifyKey, info[256:], info[:256])
	if err != nil {
		return err
	}
	decryptedInfo, err := userlib.PKEDec(userdata.RSAPrivKey, info[256:])
	if err != nil {
		return err
	}
	ln := len(decryptedInfo)
	var KeyInfo UnlockInfo
	json.Unmarshal(decryptedInfo[:ln - 32], &KeyInfo.DSKey)
	for _, v := range sharedFileList {
		if KeyInfo.DSKey == v.DSKey {
			return errors.New("Replay attack!")
		}
	}
	KeyInfo.MACKey = decryptedInfo[ln - 32:ln - 16]
	KeyInfo.DecryptionKey = decryptedInfo[ln - 16:]
	sharedFileList[filename] = KeyInfo
	byteData, _ := json.Marshal(sharedFileList)
	encryptAndStore(byteData, userdata.SharedFileListEnc)
	return nil
}

// Removes target user's access.
func (userdata *User) RevokeFile(filename string, target_username string) (err error) {
	var fileList map[string]OwnedFileMetaDataInfo
	fD, err1 := fetchAndDecrypt(userdata.FileListEnc)
	if err1 != nil {
		return err1
	}
	json.Unmarshal(fD, &fileList)
	if ownedInfo, ok := fileList[filename]; ok {
		if _, ok2 := ownedInfo.SharedUnlocks[target_username]; ok2 {
			loadedData, err := userdata.LoadFile(filename)
			if err != nil {
				return err
			}
			delete(ownedInfo.SharedUnlocks, target_username)
			getMetaData(ownedInfo.OriginalInfo, true)
			newMetaDataInfo := createSharedMDInfo()
			ownedInfo.OriginalInfo = newMetaDataInfo
			fileList[filename] = ownedInfo
			byteData, _ := json.Marshal(fileList)
			encryptAndStore(byteData, userdata.FileListEnc)
			userdata.StoreFile(filename, loadedData)
			for _, v := range ownedInfo.SharedUnlocks {
				storeMetaDataInfo(newMetaDataInfo, v)
			}
		} else {
			return errors.New("File hasn't been shared with this user!")
		}
	} else {
		return errors.New("User isn't file owner!")
	}
	return
}
