class Neopoll < Formula
    desc     "CLI poll application"
    homepage "https://github.com/neosh11/homebrew-neopoll-cli"
    url      "https://github.com/neosh11/homebrew-neopoll-cli/archive/refs/tags/v0.1.4.tar.gz"
    sha256   "97ee6213918916fce768eeb14a5904e203be261cbcb2a093672ea6f0f5e709c6"
    license  "MIT"
    depends_on "go" => :build
  
    def install
      system "go", "build", "-ldflags", "-s -w", "-o", bin/"neopoll", "."
  
      # completions
      bash_output = Utils.safe_popen_read(bin/"neopoll", "completion", "bash")
      (bash_completion/"neopoll").write bash_output
  
      zsh_output = Utils.safe_popen_read(bin/"neopoll", "completion", "zsh")
      (zsh_completion/"_neopoll").write zsh_output
  
      fish_output = Utils.safe_popen_read(bin/"neopoll", "completion", "fish")
      (fish_completion/"neopoll.fish").write fish_output
    end
  
    test do
      output = shell_output("#{bin}/neopoll --help")
      assert_match "neopoll is a CLI poll application", output
    end
  end
  