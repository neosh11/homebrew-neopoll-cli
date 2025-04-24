class Neopoll < Formula
    desc     "CLI poll application"
    homepage "https://github.com/neosh11/neopoll-cli"
    url      "https://github.com/neosh11/neopoll-cli/archive/refs/tags/v0.1.1.tar.gz"
    sha256   "5b510d5453729ed1121892f4d7f6ab6bb91915015bdb8ec29e82d77689f80700"
    license  "MIT"
    depends_on "go" => :build
  
    def install
      system "go", "build", "-ldflags", "-s -w",
             "-o", bin/"neopoll", "."
    end
  
    test do
      # simply verify it runs and prints help
      output = shell_output("#{bin}/neopoll --help")
      assert_match "neopoll is a CLI poll application", output
    end
  end
  