package proj2

// You MUST NOT change what you import.  If you add ANY additional
// imports it will break the autograder, and we will be Very Upset.

import (
	_ "encoding/hex"
	_ "encoding/json"
	_ "errors"
	"reflect"
	_ "strconv"
	_ "strings"
	"testing"

	"github.com/cs161-staff/userlib"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
)

func clear() {
	// Wipes the storage so one test does not affect another
	userlib.DatastoreClear()
	userlib.KeystoreClear()
}

func TestInit(t *testing.T) {
	clear()
	t.Log("Initialization test")

	// You can set this to false!
	userlib.SetDebugStatus(false)

	u, err := InitUser("alice", "fubar")
	if err != nil {
		// t.Error says the test fails
		t.Error("Failed to initialize user", err)
		return
	}
	// t.Log() only produces output if you run with "go test -v"
	t.Log("Got user", u)
	v, err2 := GetUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to get user", err2)
	}
	if v.Username != "alice" {
		t.Error("Incorrectly loaded user")
	}

	_, err3 := GetUser("alice", "fubar1")
	if err3 == nil {
		t.Error("Authenticated a wrong password")
	}

	_, errx1 := InitUser("alice", "fubar")
	_, errx2 := InitUser("alice", "notfubar")
	if errx1 == nil || errx2 == nil {
		t.Error("same username err")
	}

	t.Log("Got user", v)

	_, err = InitUser("__://||123", "(scm () () ())")
	_, err = InitUser("__//||123", "fubar")
	if err != nil {
		t.Error("User with weird symbols in username not allowed.")
		return
	}
	_, err = GetUser("__//||123", "(scm () () ())")
	if err == nil {
		t.Error("Authenticated user with wrong username")
		return
	}
	_, err = GetUser("__://||123", "(scm () () ())")
	if err != nil {
		t.Error("Failed to authenticate user with weird symbols in auth info.")
		return
	}

	// If you want to comment the line above,
	// write _ = u here to make the compiler happy
	// You probably want many more tests here.
}

func TestStorage(t *testing.T) {
	clear()
	u, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}

	v := []byte("This is a test")
	u.StoreFile("file1", v)

	v2, err2 := u.LoadFile("file1")
	if err2 != nil {
		t.Error("Failed to upload and download", err2)
		return
	}
	if !reflect.DeepEqual(v, v2) {
		t.Error("Downloaded file is not the same", v, v2)
		return
	}

	longData := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi quis magna sed dui dictum porta eu at odio. Fusce vehicula aliquet est, quis egestas eros dapibus vestibulum. Aliquam sapien augue, pulvinar vitae porta sodales, feugiat vel lectus. Phasellus tristique, risus quis sollicitudin maximus, diam risus efficitur massa, porta cursus nibh ante at leo. Quisque nec velit eget mauris egestas lacinia in ac lorem. Morbi mollis tempus ante. Aenean fermentum nisl a mi maximus lobortis. Proin sagittis purus enim, vel cursus arcu porttitor blandit. In nulla tellus, euismod quis sapien et, interdum ornare ex. Etiam malesuada augue id ipsum pharetra, ac volutpat purus scelerisque. Curabitur ac enim ac metus condimentum faucibus sit amet ac magna. Vivamus luctus laoreet orci, eget ornare odio viverra id. Suspendisse mattis tortor vel lacus volutpat egestas. Suspendisse ac neque velit. Mauris venenatis lectus a porttitor tincidunt. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Mauris venenatis tempus lectus eleifend commodo. Curabitur urna velit, consectetur ac ornare in, rhoncus a nunc. Maecenas congue nibh ut dolor malesuada molestie. Donec at erat quis metus volutpat euismod. In hac habitasse platea dictumst. Etiam sit amet quam sagittis, pulvinar sem quis, tempus elit. Curabitur vel elit sagittis, eleifend purus nec, maximus arcu. Pellentesque imperdiet et leo at venenatis. Curabitur pretium, justo fringilla suscipit laoreet, nisl eros consequat est, tempor tempor elit sapien laoreet erat. Integer volutpat rutrum tempor. Curabitur ut interdum tellus. Praesent interdum nibh nec justo vulputate, in rutrum purus aliquet. Fusce rutrum mauris eu blandit feugiat. Curabitur tincidunt fringilla elit, vitae malesuada lacus laoreet ac. Aliquam tristique sem lacus, sed consectetur leo convallis eu. Sed sed justo sed libero scelerisque sodales. Cras non orci malesuada, ultrices eros nec, semper ipsum. Duis ut auctor elit. Proin mollis erat justo, id suscipit erat vulputate sed. Nunc et ante odio. Praesent porttitor, magna et lobortis porttitor, nisl nisl facilisis mauris, ac suscipit orci velit viverra risus. Donec congue felis ac ex finibus posuere sit amet quis neque. Vivamus congue iaculis velit sed rhoncus. Phasellus et massa magna. Etiam eget mi vel augue imperdiet rutrum vitae et eros. Pellentesque eget sem lobortis, lobortis nibh sed, facilisis risus. Nullam sed ornare lectus. Nullam sollicitudin purus nec diam pretium, sed imperdiet tortor semper. Phasellus non molestie risus. Quisque iaculis, tortor in semper ullamcorper, diam mauris tristique erat, eget gravida risus ligula semper ante. Nulla facilisi. Mauris fermentum pellentesque posuere. Proin faucibus tempor elit, eget placerat eros convallis eget. Pellentesque a eros et augue dapibus elementum. Pellentesque venenatis sollicitudin lorem, ac posuere tortor viverra nec. Aliquam erat volutpat. Integer egestas ante feugiat malesuada malesuada. Donec dignissim ligula risus, nec euismod mauris congue vel. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam faucibus pharetra egestas. Mauris et blandit risus, nec eleifend quam. Ut blandit semper nisi vitae mattis.")
	finalData := longData[:]
	for x := 0; x < 1000; x += 1 {
		finalData = append(finalData[:], longData[:]...)
	}
	t.Log(len(finalData))
	u.StoreFile("file2", finalData)
	dt, err3 := u.LoadFile("file2")
	t.Log(len(dt))
	if !reflect.DeepEqual(dt, finalData) {
		t.Error("Long data store does not work.")
	}
	if err3 != nil {
		t.Error("Load data does not work for long data")
	}
}

func TestInvalidFile(t *testing.T) {
	clear()
	u, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}

	_, err2 := u.LoadFile("this file does not exist")
	if err2 == nil {
		t.Error("Downloaded a ninexistent file", err2)
		return
	}
}

func TestShare(t *testing.T) {
	clear()
	u, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}
	u2, err2 := InitUser("bob", "foobar")
	if err2 != nil {
		t.Error("Failed to initialize bob", err2)
		return
	}

	v := []byte("This is a test")
	u.StoreFile("file1", v)

	var v2 []byte
	var magic_string string

	v, err = u.LoadFile("file1")
	if err != nil {
		t.Error("Failed to download the file from alice", err)
		return
	}

	magic_string, err = u.ShareFile("file1", "bob")
	if err != nil {
		t.Error("Failed to share the a file", err)
		return
	}
	err = u2.ReceiveFile("file2", "alice", magic_string)
	if err != nil {
		t.Error("Failed to receive the share message", err)
		return
	}

	v2, err = u2.LoadFile("file2")
	if err != nil {
		t.Error("Failed to download the file after sharing", err)
		return
	}
	if !reflect.DeepEqual(v, v2) {
		t.Error("Shared file is not the same", v, v2)
		return
	}

	// aidan's tests below
	clear()
	// alice sharing file1 with bob
	alice, err := InitUser("alice", "fubar")
	bob, err := InitUser("bob", "fubar")
	alice.StoreFile("file1", []byte("this is file1."))
	magic_string, err = alice.ShareFile("file1", "bob")
	if err != nil {
		t.Error("failed to share file1", err)
		return
	}
	err = bob.ReceiveFile("file1", "alice", magic_string)
	if err != nil {
		t.Error("failed to receive file1", err)
		return
	}

	// alice sharing file1 with bob a second time, bob receives with same name
	magic_string, err = alice.ShareFile("file1", "bob")
	if err != nil {
		t.Error("failed to share file1 the second time", err)
		return
	}
	err = bob.ReceiveFile("file1", "alice", magic_string)
	if err == nil { // error expected
		t.Error("received file a second time using the same name", err)
		return
	}
	// bob receives file1 under a different name (using token that previously errored out)
	err = bob.ReceiveFile("file1_copy", "alice", magic_string)
	if err != nil {
		t.Error("failed to receive file1 a second time, this time called 'file1_copy'", err)
		return
	}
	v1, err := bob.LoadFile("file1")
	if err != nil {
		t.Error("Failed to load file1", err)
		return
	}
	v2, err = bob.LoadFile("file1_copy")
	if err != nil {
		t.Error("Failed to load file1_copy", err)
		return
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Error("file1 and file1_copy shuold be the same", v1, v2)
		return
	}
	v3, err := alice.LoadFile("file1")
	if !reflect.DeepEqual(v1, v3) {
		t.Error("Alice's file1 should be the same as Bob's file1", v1, v3)
		return
	}
	if !reflect.DeepEqual(v2, v3) {
		t.Error("Alice's file1 should be the same as Bob's file1_copy", v2, v3)
		return
	}

	// test receiving a file with a name that is already used
	alice.StoreFile("file2", []byte("this is file2"))
	magic_string, err = alice.ShareFile("file2", "bob")
	if err != nil {
		t.Error("failed to share file2", err)
		return
	}
	err = bob.ReceiveFile("file1", "alice", magic_string)
	if err == nil { // error expected
		t.Error("using the same name for a received file did not error", err)
		return
	}

	// test that bob can share file owned by alice
	charlie, err := InitUser("charlie", "fubar")
	magic_string, err = bob.ShareFile("file1", "charlie")
	if err != nil {
		t.Error("failed to share file1", err)
		return
	}
	err = charlie.ReceiveFile("file1", "bob", magic_string)
	if err != nil {
		t.Error("charlie failed to receive file1 from bob", err)
		return
	}
	v1, err = bob.LoadFile("file1")
	v2, err = charlie.LoadFile("file1")
	if !reflect.DeepEqual(v1, v2) {
		t.Error("charlie and bob's file1's shuold be the same", v1, v2)
		return
	}

	// test that changes made by A, B, C can be seen by everyone
	err = alice.AppendFile("file1", []byte("alice appended this."))
	v1, err = alice.LoadFile("file1")
	v2, err = bob.LoadFile("file1")
	v3, err = charlie.LoadFile("file1")
	check := append([]byte("this is file1."), []byte("alice appended this.")...)
	if !reflect.DeepEqual(v1, check) {
		t.Error("alice's append failed")
		return
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Error("alice's changes not reflected across all shared users")
		return
	}
	if !reflect.DeepEqual(v2, v3) {
		t.Error("alice's changes not reflected across all shared users")
		return
	}
	err = charlie.AppendFile("file1", []byte("charlie appended this."))
	v1, err = charlie.LoadFile("file1")
	v2, err = bob.LoadFile("file1_copy")
	if !reflect.DeepEqual(v1, v2) {
		t.Error("charlie's changes not reflected in bob's file1_copy")
		return
	}

	// test charlie storing to (overwriting) an existing shared file
	check = []byte("charlie overwriting file1.")
	charlie.StoreFile("file1", check)
	v1, err = charlie.LoadFile("file1")
	if !reflect.DeepEqual(v1, check) {
		t.Error("charlie failed to overwrite file1 using StoreFile")
		return
	}
	v2, err = alice.LoadFile("file1")
	if !reflect.DeepEqual(v1, v2) {
		t.Error("charlie's newly stored 'file1' is not reflected in alice's file1")
		return
	}

	// test receiving with an incorrect sender
	magic_string, err = alice.ShareFile("file2", "charlie")
	err = charlie.ReceiveFile("file2", "bob", magic_string)
	if err == nil {
		t.Error("charlie received file2 from alice even though bob was passed in as the sender")
		return
	}
	v1, err = charlie.LoadFile("file2")
	if err == nil {
		t.Error("charlie was able to load a file he doesn't have")
		return
	}

	// test receiving with an incorrect magic_string
	charlie.StoreFile("file3", []byte("this is file3."))
	magic_string, err = charlie.ShareFile("file3", "alice")
	fake_magic := string(append([]byte(magic_string)[:], []byte("1")...))
	err = alice.ReceiveFile("file3", "charlie", fake_magic)
	if err == nil {
		t.Error("alice was able to receive a file with an incorrect magic string")
		return
	}

	//
	// testing with multiple instances
	//

	// 2 instances of alice load, append, and overwrite a file
	check = []byte("this is file4.")
	alice.StoreFile("file4", check)
	alice1, err := GetUser("alice", "fubar")
	if err != nil {
		t.Error("first call to GetUser for alice failed")
		return
	}
	alice2, err := GetUser("alice", "fubar")
	if err != nil {
		t.Error("second call to GetUser for alice failed")
		return
	}
	v1, err = alice1.LoadFile("file4")
	v2, err = alice2.LoadFile("file4")
	if !reflect.DeepEqual(v1, check) {
		t.Error("alice1's loaded file4 is incorrect")
		return
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Error("alice1's and alice2's file4s are different")
		return
	}

	// test with appending
	check = append(check[:], []byte("alice1 appended this.")...)
	_ = alice1.AppendFile("file4", []byte("alice1 appended this."))
	v1, err = alice1.LoadFile("file4")
	v2, err = alice2.LoadFile("file4")
	if !reflect.DeepEqual(v2, check) {
		t.Error("alice2's loaded file4 is incorrect after alice1 appended")
		return
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Error("alice1's and alice2's file4s are different after appending")
		return
	}

	// alice2 overwrites file4 by calling store
	check = []byte("hi")
	alice2.StoreFile("file4", check)
	v1, err = alice1.LoadFile("file4")
	if !reflect.DeepEqual(v1, check) {
		t.Error("alice1's file4 was not overwritten by alice2", check, v1)
		return
	}
}

func TestAppend(t *testing.T) {
	clear()
	_, err1 := InitUser("a", "hello777$$$")
	a, err2 := GetUser("a", "hello777$$$")
	data1 := []byte("This is part a. ")
	data2 := []byte("This is part b. ")
	a.StoreFile("one", data1)
	a.AppendFile("one", data2)
	check, err3 := a.LoadFile("one")
	if !reflect.DeepEqual(check, append(data1[:], data2[:]...)) {
		t.Error("Appended File is not the same", check, append(data1[:], data2[:]...))
	}
	if err1 != nil || err2 != nil || err3 != nil {
		t.Error("Something else failed")
	}
	b, err4 := InitUser("b", "goodbye123^^^")
	longData := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi quis magna sed dui dictum porta eu at odio. Fusce vehicula aliquet est, quis egestas eros dapibus vestibulum. Aliquam sapien augue, pulvinar vitae porta sodales, feugiat vel lectus. Phasellus tristique, risus quis sollicitudin maximus, diam risus efficitur massa, porta cursus nibh ante at leo. Quisque nec velit eget mauris egestas lacinia in ac lorem. Morbi mollis tempus ante. Aenean fermentum nisl a mi maximus lobortis. Proin sagittis purus enim, vel cursus arcu porttitor blandit. In nulla tellus, euismod quis sapien et, interdum ornare ex. Etiam malesuada augue id ipsum pharetra, ac volutpat purus scelerisque. Curabitur ac enim ac metus condimentum faucibus sit amet ac magna. Vivamus luctus laoreet orci, eget ornare odio viverra id. Suspendisse mattis tortor vel lacus volutpat egestas. Suspendisse ac neque velit. Mauris venenatis lectus a porttitor tincidunt. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Mauris venenatis tempus lectus eleifend commodo. Curabitur urna velit, consectetur ac ornare in, rhoncus a nunc. Maecenas congue nibh ut dolor malesuada molestie. Donec at erat quis metus volutpat euismod. In hac habitasse platea dictumst. Etiam sit amet quam sagittis, pulvinar sem quis, tempus elit. Curabitur vel elit sagittis, eleifend purus nec, maximus arcu. Pellentesque imperdiet et leo at venenatis. Curabitur pretium, justo fringilla suscipit laoreet, nisl eros consequat est, tempor tempor elit sapien laoreet erat. Integer volutpat rutrum tempor. Curabitur ut interdum tellus. Praesent interdum nibh nec justo vulputate, in rutrum purus aliquet. Fusce rutrum mauris eu blandit feugiat. Curabitur tincidunt fringilla elit, vitae malesuada lacus laoreet ac. Aliquam tristique sem lacus, sed consectetur leo convallis eu. Sed sed justo sed libero scelerisque sodales. Cras non orci malesuada, ultrices eros nec, semper ipsum. Duis ut auctor elit. Proin mollis erat justo, id suscipit erat vulputate sed. Nunc et ante odio. Praesent porttitor, magna et lobortis porttitor, nisl nisl facilisis mauris, ac suscipit orci velit viverra risus. Donec congue felis ac ex finibus posuere sit amet quis neque. Vivamus congue iaculis velit sed rhoncus. Phasellus et massa magna. Etiam eget mi vel augue imperdiet rutrum vitae et eros. Pellentesque eget sem lobortis, lobortis nibh sed, facilisis risus. Nullam sed ornare lectus. Nullam sollicitudin purus nec diam pretium, sed imperdiet tortor semper. Phasellus non molestie risus. Quisque iaculis, tortor in semper ullamcorper, diam mauris tristique erat, eget gravida risus ligula semper ante. Nulla facilisi. Mauris fermentum pellentesque posuere. Proin faucibus tempor elit, eget placerat eros convallis eget. Pellentesque a eros et augue dapibus elementum. Pellentesque venenatis sollicitudin lorem, ac posuere tortor viverra nec. Aliquam erat volutpat. Integer egestas ante feugiat malesuada malesuada. Donec dignissim ligula risus, nec euismod mauris congue vel. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam faucibus pharetra egestas. Mauris et blandit risus, nec eleifend quam. Ut blandit semper nisi vitae mattis.")
	finalData := longData[:]
	for x := 0; x < 1000; x += 1 {
		finalData = append(finalData[:], longData[:]...)
		t.Log(x)
	}
	t.Log("STart storing")
	b.StoreFile("txt", finalData)
	t.Log("Finished Storing")
	b.AppendFile("txt", data1)
	t.Log("Finished Appending")
	check2, err5 := b.LoadFile("txt")
	newDt := append(finalData[:], data1[:]...)
	if !reflect.DeepEqual(check2, newDt) {
		t.Error("Appended File is not the same", len(check), len(newDt))
	}
	b.AppendFile("txt", finalData)
	newDt = append(newDt[:], finalData[:]...)
	check3, err6 := b.LoadFile("txt")
	if !reflect.DeepEqual(check3, newDt) {
		t.Error("Appended File is not the same", check3, newDt)
	}
	if err4 != nil || err5 != nil || err6 != nil {
		t.Error("Something else failed")
	}
	err7 := b.AppendFile("txt1", finalData)
	if err7 == nil {
		t.Error("Appending to a nonexistant file!")
	}
}

func TestRevoke(t *testing.T) {
	clear()
	t.Log("Revoke test")
	u, err := InitUser("alice", "1234567")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}
	u2, err2 := InitUser("bob", "9876543")
	if err2 != nil {
		t.Error("Failed to initialize user", err)
		return
	}
	w := []byte("This is a test fdjsfidsofjs")
	u.StoreFile("file1", w)

	var v2 []byte
	var magic_string string
	magic_string, err = u.ShareFile("file1", "bob")
	if err != nil {
		t.Error("Failed to share the a file", err)
		return
	}
	err = u2.ReceiveFile("file2", "alice", magic_string)
	if err != nil {
		t.Error("Failed to receive the share message", err)
		return
	}

	v2, err = u2.LoadFile("file2")
	if err != nil {
		t.Error("Failed to download the file after sharing", err)
		return
	}
	if !reflect.DeepEqual(w, v2) {
		t.Error("Shared file is not the same", w, v2)
		return
	}
	u3, _ := InitUser("dave", "dave123")
	ms, _ := u2.ShareFile("file2", "dave")
	_ = u3.ReceiveFile("file3", "bob", ms)
	u4, _ := InitUser("george", "george456")
	ms2, _ := u.ShareFile("file1", "george")
	err26 := u4.ReceiveFile("file4", "alice", ms2)
	if err26 != nil {
		t.Error("Tree share not right", err26)
	}
	v3, err17 := u3.LoadFile("file3")
	if err17 != nil {
		t.Error("Failed to download the file after sharing", err17)
		return
	}
	if !reflect.DeepEqual(w, v3) {
		t.Error("Shared file is not the same", w, v3)
		return
	}
	newData := []byte("This is a differenr file !!!!!!!")
	u.RevokeFile("file1", "bob")
	v7, err89 := u.LoadFile("file1")
	if err89 != nil {
		t.Error("Failed to download the file after sharing..", err89)
		return
	}
	if !reflect.DeepEqual(w, v7) {
		t.Error("Shared file is not the same", w, v7)
		return
	}
	u.StoreFile("file1", newData)
	file_check, err9 := u2.LoadFile("file2")
	if err9 == nil && reflect.DeepEqual(file_check, newData) {
		t.Error("Revoke did not work properly", string(file_check))
		return
	}
	file_check2, err99 := u3.LoadFile("file3")
	if err99 == nil && reflect.DeepEqual(file_check2, newData) {
		t.Error("Revoke did not work properly", string(file_check2))
		return
	}

	file_check5, err999 := u4.LoadFile("file4")
	if err999 != nil || !reflect.DeepEqual(file_check5, newData) {
		t.Error("Revoke did not work properly", err999)
		return
	}

	// aidan's tests (users named aidan, bill, charlie, david, eve, fish)
	clear()
	// test basic tree structure
	// 	    D
	// 	   /
	// 	  B -- E -- F
	//   /
	// A
	//   \
	//    C
	aidan, _ := InitUser("aidan", "fubar")
	bill, _ := InitUser("bill", "fubar")
	charlie, _ := InitUser("charlie", "fubar")
	david, _ := InitUser("david", "fubar")
	eve, _ := InitUser("eve", "fubar")
	fish, _ := InitUser("fish", "fubar")

	aidan.StoreFile("file1", []byte("this is file1.")) // aidan creates file1

	magic_string, err = aidan.ShareFile("file1", "bill") // aidan shares with bill
	bill.ReceiveFile("file1", "aidan", magic_string)
	magic_string, err = aidan.ShareFile("file1", "charlie") // aidan shares with charlie
	charlie.ReceiveFile("file1", "aidan", magic_string)
	magic_string, err = bill.ShareFile("file1", "david") // bill shares with david
	david.ReceiveFile("file1", "bill", magic_string)
	magic_string, err = bill.ShareFile("file1", "eve") // bill shares with eve
	eve.ReceiveFile("file1", "bill", magic_string)
	magic_string, err = eve.ShareFile("file1", "fish") // eve shares with fish
	fish.ReceiveFile("file1", "eve", magic_string)

	v1, err := aidan.LoadFile("file1")
	if err != nil {
		t.Error("Failed to download the file after sharing..")
		return
	}
	v2, err = eve.LoadFile("file1")
	if err != nil {
		t.Error("Failed to download the file after sharing..")
		return
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Error("Shared file is not the same")
		return
	}
	check := append(v1[:], []byte("fish appended this.")...)
	err = fish.AppendFile("file1", []byte("fish appended this."))
	if err != nil {
		t.Error("Recipient failed to append to file")
		return
	}
	v1, err = fish.LoadFile("file1")
	if !reflect.DeepEqual(v1, check) {
		t.Error("Modifications by recipient were not saved")
		return
	}
	v2, err = david.LoadFile("file1")
	if !reflect.DeepEqual(v1, v2) {
		t.Error("Modifications by one recipient are not seen by other recipients")
		return
	}
	v3, err = aidan.LoadFile("file1")
	if !reflect.DeepEqual(v1, v3) {
		t.Error("Modifications by one recipient are not propogated up to owner of file")
		return
	}

	// test non-owner calling revoke
	err = eve.RevokeFile("file1", "fish")
	if err == nil { // error expected
		t.Error("No error was thrown when non-owner of a file tries to revoke permissions.")
		return
	}
	v1, err = fish.LoadFile("file1") // fish should still be able to access file
	if err != nil {
		t.Error("Non-owner successfully revoked another recipient's acccess")
		return
	}

	// now aidan will revoke bill -> expect david, eve, and fish to also lose access
	err = aidan.RevokeFile("file1", "bill")
	if err != nil {
		t.Error("Owner failed to revoke permissions")
		return
	}

	check, _ = aidan.LoadFile("file1") // will use this below

	// check that bill and one of his children (eve) can't append to file1 anymore
	err1 := bill.AppendFile("file1", []byte("bill tried to append this after being revoked."))
	err2 = eve.AppendFile("file1", []byte("eve tried to append this after being revoked."))
	if err1 == nil || err2 == nil {
		t.Error("Revoked users could still append to file")
		return
	}

	// all users load file 1
	aidan_file1, _ := aidan.LoadFile("file1")
	_, bill_err := bill.LoadFile("file1")
	charlie_file1, _ := charlie.LoadFile("file1")
	_, david_err := david.LoadFile("file1")
	_, eve_err := eve.LoadFile("file1")
	_, fish_err := fish.LoadFile("file1")

	// check that bill and none of his children can load the file
	if bill_err == nil || david_err == nil || eve_err == nil || fish_err == nil {
		t.Error("At least one of these revoked user's could still load the file")
		return
	}

	// check that bill can no longer share the file that he was revoked from
	magic_string, err = bill.ShareFile("file1", "fish")
	if err == nil {
		t.Error("Revoked user can still share the file")
		return
	}

	// check that even if a revoked user calls StoreFile with the revoked filename, nothing happens
	david.StoreFile("file1", []byte("david tried to storefile with the revoked filename"))
	v1, err = david.LoadFile("file1")
	if err == nil {
		t.Error("Revoked user can still call StoreFile() with the same filename as the revoked file", v1)
		return
	}

	// check that aidan and charlie's file1 are unaffected by the revoke
	if !reflect.DeepEqual(aidan_file1, charlie_file1) {
		t.Error("aidan and charlie's file's are different after revoking bill")
		return
	}

	// check that aidan and charlie can still use file1 as usual
	check = append(check[:], []byte("aidan appended this after revoking bill.")...)
	err = aidan.AppendFile("file1", []byte("aidan appended this after revoking bill."))
	v1, err = charlie.LoadFile("file1")

	if !reflect.DeepEqual(v1, check) {
		t.Error("aidan and charlie can't append like normal after revoking bill")
		return
	}

	// test resharing file 1 with bill
	magic_string, err = aidan.ShareFile("file1", "bill")
	err = bill.ReceiveFile("file1_reshare", "aidan", magic_string)
	if err != nil {
		t.Error("Bill failed to receive the same file again under the same name")
		return
	}

	v1, err = aidan.LoadFile("file1")
	v2, err = bill.LoadFile("file1_reshare")
	if !reflect.DeepEqual(v1, v2) {
		t.Error("aidan and bill have different file1's after resharing")
		return
	}

	// bill shares again with one of his children (fish)
	magic_string, err = bill.ShareFile("file1_reshare", "fish")
	err = fish.ReceiveFile("file1", "bill", magic_string) // expected to fail
	if err == nil {
		t.Error("fish was able to receive the file as 'file1' again")
		return
	}
	err = fish.ReceiveFile("file1_reshare", "bill", magic_string)
	if err != nil {
		t.Error("fish was unable to receive the reshared file1")
		return
	}

}

func TestAttacksKunal(t *testing.T) {
	clear()
	a, _ := InitUser("alice", "123")
	b, _ := InitUser("bob", "123")
	f1 := []byte("fdsf")
	f2 := []byte("")
	a.StoreFile("", f1)
	a.StoreFile("1", f2)
	q1, check1 := a.LoadFile("")
	q2, check2 := a.LoadFile("1")
	if check1 != nil || check2 != nil {
		t.Error("empty file / name")
	}
	if !reflect.DeepEqual(string(q1), string(f1)) || !reflect.DeepEqual(string(q2), string(f2)) {
		t.Error("error empty", q1, f1, q2, f2)
	}
	a.StoreFile("f", []byte("12345"))
	_, er1 := a.LoadFile("fdsf")
	er2 := a.AppendFile("fdsfds", []byte("fdsfaf"))
	if er1 == nil || er2 == nil {
		t.Error("Append/Loading to nonexistant file")
	}
	_, e := a.ShareFile("g", "bob")
	_, e1 := a.ShareFile("g", "charlie")
	if e == nil || e1 == nil {
		t.Error("sharing wrong")
	}
	s, _ := a.ShareFile("f", "bob")
	InitUser("pp", "xyz")
	e424 := b.ReceiveFile("f", "pp", s)
	if e424 == nil {
		t.Error("Authenticity fail")
	}
	e4 := b.ReceiveFile("f", "hola", s)
	s1 := s + "4"
	e22 := b.ReceiveFile("f", "alice", s1)
	if e22 == nil || e4 == nil {
		t.Error("Bad receive")
	}
	b.ReceiveFile("f", "alice", s)
	e = a.RevokeFile("f", "cool")
	if e == nil {
		t.Error("NE user")
	}
	a.RevokeFile("f", "bob")
	e = a.RevokeFile("f", "bob")
	if e == nil {
		t.Error("Revoked an already revoked file :(")
	}
	b.AppendFile("f", []byte("456"))
	l, _ := b.LoadFile("f")
	p, _ := a.LoadFile("f")
	t.Log(l, p)
	dummy := []byte("fdsf")
	for k, v := range userlib.DatastoreGetMap() {
		userlib.DatastoreSet(k, append(v[:], dummy[:]...))
	}
	_, e100 := a.LoadFile("f")
	_, e101 := b.LoadFile("f")
	e102 := a.AppendFile("f", dummy)
	e103 := b.AppendFile("f", dummy)
	_, e104 := GetUser("alice", "123")
	t.Log(e100, e101, e102, e103, e104)
	if e100 == nil || e101 == nil || e102 == nil || e103 == nil || e104 == nil {
		t.Error("Something went wrong")
	}
	fileH := []byte("File H")
	a.StoreFile("h", fileH)
	e = a.RevokeFile("h", "bob")
	e1 = a.RevokeFile("i", "bob")
	if e == nil || e1 == nil {
		t.Error("Revoked an unshared/nonexistant file")
	}
}

func TestMultipleInstancesKunal(t *testing.T) {
	clear()
	u1, e1 := InitUser("alice", "hola")
	u2, e2 := GetUser("alice", "hola")
	u3, e3 := GetUser("alice", "hola")
	if e1 != nil || e2 != nil || e3 != nil {
		t.Error("Init/Get err")
	}

	f1 := []byte("hello")
	u1.StoreFile("f", f1)
	r1, e4 := u1.LoadFile("f")
	r2, e5 := u2.LoadFile("f")
	r3, e6 := u3.LoadFile("f")
	if e4 != nil || e5 != nil || e6 != nil {
		t.Error("Load/Store err")
	}
	if !reflect.DeepEqual(string(r1), string(f1)) {
		t.Error("Mismatch u1")
	}
	if !reflect.DeepEqual(string(r2), string(f1)) {
		t.Error("Mismatch u2")
	}
	if !reflect.DeepEqual(string(r3), string(f1)) {
		t.Error("Mismatch u3")
	}

	app := []byte(" world")
	f2 := []byte("hello world")
	e7 := u2.AppendFile("f", app)
	r4, e8 := u1.LoadFile("f")
	r5, e9 := u2.LoadFile("f")
	r6, e0 := u3.LoadFile("f")
	if e7 != nil || e8 != nil || e9 != nil || e0 != nil {
		t.Error("Load/Store/Append err")
	}
	if !reflect.DeepEqual(string(r4), string(f2)) {
		t.Error("Mismatch u1")
	}
	if !reflect.DeepEqual(string(r5), string(f2)) {
		t.Error("Mismatch u2")
	}
	if !reflect.DeepEqual(string(r6), string(f2)) {
		t.Error("Mismatch u3")
	}
}
func TestDSAttacksKunal(t *testing.T) {
	clear()
	a, _ := InitUser("a", "1")
	b, _ := InitUser("b", "2")
	m1 := userlib.DatastoreGetMap()
	attack := []byte("x")
	for k, v := range m1 {
		userlib.DatastoreSet(k, append(v[:], attack[:]...))
	}
	_, e1 := GetUser("a", "1")
	_, e2 := GetUser("b", "2")
	if e1 == nil || e2 == nil {
		t.Error("Successful attack!")
	}

	for k, v := range m1 {
		userlib.DatastoreSet(k, v)
	}
	f1 := []byte("file 1")
	f2 := []byte("file 2")
	a.StoreFile("f1", f1)
	b.StoreFile("f1", f2)
	m2 := userlib.DatastoreGetMap()
	for k, v := range m2 {
		userlib.DatastoreSet(k, append(v[:], attack[:]...))
	}
	for k, v := range m1 {
		userlib.DatastoreSet(k, v)
	}
	add := []byte("additional")
	e1 = a.AppendFile("f1", add)
	e2 = a.AppendFile("f1", add)
	_, e3 := a.LoadFile("f1")
	_, e4 := b.LoadFile("f1")
	if e1 == nil || e2 == nil || e3 == nil || e4 == nil {
		t.Error("Successful attack!")
	}
}

func TestDSAttacksAidan(t *testing.T) {
	// delete either a file block/metadata from the datastore
	clear()
	alice, _ := InitUser("alice", "fubar")
	bob, _ := InitUser("bob", "fubar")
	charlie, _ := InitUser("charlie", "fubar")

	m1 := userlib.DatastoreGetMap()
	userInfo := make(map[uuid.UUID][]byte) // map that contains DS content after users are created; will use this to isolate file_info on the DS
	for k, v := range m1 {
		userInfo[k] = v
	}

	alice.StoreFile("file1", []byte("this is file1.")) // alice creates file1

	m2 := userlib.DatastoreGetMap()
	fileInfo := make(map[uuid.UUID][]byte) // map that contains only DS content for file1
	for k, v := range m2 {
		_, ok := userInfo[k]
		if !ok {
			fileInfo[k] = v
		}
	}

	magicString, _ := alice.ShareFile("file1", "alice") // alice shares file1 with bob
	_ = bob.ReceiveFile("file1", "alice", magicString)

	// select a key from fileInfo and remove it
	counter := 0
	var removeKey uuid.UUID
	var removeVal []byte
	for k, v := range fileInfo {
		if counter == 0 {
			removeKey = k
			removeVal = v
		}
	}
	t.Log(len(removeKey))
	userlib.DatastoreDelete(removeKey)

	// now bob will try all the file operations
	v1, err := bob.LoadFile("file1")
	if err == nil {
		t.Error("Bob was able to load file without all the file info", v1)
		return
	}

	err = bob.AppendFile("file1", []byte("bob tried to append this."))
	if err == nil {
		t.Error("Bob was able to append to file1 without all the file info")
		return
	}

	magicString, err1 := bob.ShareFile("file1", "charlie")
	err2 := charlie.ReceiveFile("file1", "bob", magicString)

	if err1 == nil || err2 == nil {
		t.Error("Sharing was successful even without all the file info")
		return
	}

	userlib.DatastoreSet(removeKey, removeVal) // add the removed entry back to datastore

}

func TestSendReceive(t *testing.T) {
	clear()
	a, _ := InitUser("a", "1")
	b, _ := InitUser("b", "1")
	c, _ := InitUser("c", "1")
	f1 := []byte("hola")
	f2 := []byte("aloha")
	a.StoreFile("hh", []byte(""))
	a.StoreFile("f", f1)
	a.StoreFile("g", f2)
	a.ShareFile("ii", "b")
	a.ShareFile("f", "d")
	s1, _ := a.ShareFile("f", "b")
	s2, _ := a.ShareFile("g", "b")
	e1 := c.ReceiveFile("h", "a", s1)
	e2 := b.ReceiveFile("h", "b", s1)
	e3 := b.ReceiveFile("h", "a", s2)
	e4 := b.ReceiveFile("i", "a", s2)

	if e1 == nil || e2 == nil || e3 != nil || e4 == nil {
		t.Error("S/R")
	}

}
