package pkg

import (
	"fmt"
	"k8s.io/kubernetes/pkg/util/mount"
	"os"
	"os/exec"
	"time"
	"k8s.io/klog/v2"
)


//mount interface
// Mounter interface which can be implemented
// by the different mounter types
type Mounter interface {
	Stage(stagePath string) error
	Unstage(stagePath string) error
	Mount(source string, target string) error
}

// newMounter returns a new mounter depending on the mounterType parameter
func newMounter(bucket string, secrets map[string]string) (Mounter, error) {
	return newS3fsMounter(bucket, secrets)
}

func fuseMount(path string, command string, args []string) error {
	cmd := exec.Command(command, args...)
	klog.V(4).Infof("Mounting fuse with command: %s and args: %s", command, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error fuseMount command: %s\nargs: %s\noutput: %s", command, args, out)
	}

	return waitForMount(path, 10*time.Second)
}

func fuseUnmount(path string) error {
	if err := mount.New("").Unmount(path); err != nil {
		return err
	}
	// as fuse quits immediately, we will try to wait until the process is done
	process, err := findFuseMountProcess(path)
	if err != nil {
		klog.Errorf("Error getting PID of fuse mount: %s", err)
		return nil
	}
	if process == nil {
		klog.Warningf("Unable to find PID of fuse mount %s, it must have finished already", path)
		return nil
	}
	klog.Infof("Found fuse pid %v of mount %s, checking if it still runs", process.Pid, path)
	return waitForProcess(process, 1)
}


// Implements Mounter interface
type s3fsMounter struct {
	bucket          string
	endpoint        string
	pwFileContent   string
	accessKeyID     string
	secretAccessKey string
}

const (
	CmdS3FS = "s3fs"
)

// todo update what minio secrets from
func newS3fsMounter(bucket string, secrets map[string]string) (Mounter, error) {
	return &s3fsMounter{
		bucket:          bucket,
		//accessKeyID:     secrets["accessKeyID"],
		//secretAccessKey: secrets["secretAccessKey"],
		//endpoint:        secrets["endpoint"],

		accessKeyID:     os.Getenv("MINIO_ACCESS_KEY"),
		secretAccessKey: os.Getenv("MINIO_SECRET_KEY"),
		endpoint:        os.Getenv("MINIO_ENDPOINT"),
	}, nil
}

func (s3fs *s3fsMounter) Stage(stageTarget string) error {
	return nil
}

func (s3fs *s3fsMounter) Unstage(stageTarget string) error {
	return nil
}

func (s3fs *s3fsMounter) Mount(source string, target string) error {
	if err := writes3fsPass(s3fs.accessKeyID + ":" + s3fs.secretAccessKey); err != nil {
		return err
	}
	args := []string{
		fmt.Sprintf("%s", s3fs.bucket),
		fmt.Sprintf("%s", target),
		"-o", "use_path_request_style",
		"-o", fmt.Sprintf("url=http://%s", s3fs.endpoint),
		"-o", "curldbg",
		"-o", "allow_other",
		"-o", "mp_umask=000",
	}
	return fuseMount(target, CmdS3FS, args)
}

func writes3fsPass(pwFileContent string) error {
	pwFileName := fmt.Sprintf("%s/.passwd-s3fs", os.Getenv("HOME"))
	pwFile, err := os.OpenFile(pwFileName, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	_, err = pwFile.WriteString(pwFileContent)
	if err != nil {
		return err
	}
	pwFile.Close()
	return nil
}
