class Neopoll < Formula
    desc     "CLI poll application"
    homepage "https://github.com/neosh11/neopoll-cli"
    url      "https://github.com/neosh11/neopoll-cli/archive/refs/tags/v0.1.1.tar.gz"
    sha256   "0019dfc4b32d63c1392aa264aed2253c1e0c2fb09216f8e2cc269bbfb8bb49b5"
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
  