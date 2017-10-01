package main

import (
	"github.com/pkg/errors"
	"os/exec"
)

func EncodeVideo(image, speech, output_path string) error {
	cmd_args := []string{
		"-loop", "1",
		"-i", image,
		"-i", speech,
		"-c:a", "aac",
		"-strict", "experimental",
		"-b:a", "512k",
		"-pix_fmt", "yuv420p",
		"-movflags", "+faststart",
		"-shortest", output_path,
	}
	out, err := exec.Command("ffmpeg", cmd_args...).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "Faild video encode.\n"+string(out))
	}
	return nil
}
