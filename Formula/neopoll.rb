class Neopoll < Formula
    desc     "CLI poll application"
    homepage "https://github.com/neosh11/homebrew-neopoll-cli"
    url      "https://github.com/neosh11/homebrew-neopoll-cli/archive/refs/tags/v0.1.5.tar.gz"
    sha256   "57fdb58f9f412951f20f8285ed04b14493f3d70fb68d10c56a525b1e9cd9e5ac"
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
  