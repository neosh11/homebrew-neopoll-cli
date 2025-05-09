class Neopoll < Formula
    desc     "CLI poll application"
    homepage "https://github.com/neosh11/homebrew-neopoll-cli"
    url      "https://github.com/neosh11/homebrew-neopoll-cli/archive/refs/tags/v0.1.6.tar.gz"
    sha256   "cee40d3b056b674d5854aab921a5736be02fb4fe56fdced0e38535261f6ed1b4"
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
  