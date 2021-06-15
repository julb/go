require 'digest'
require 'rake/clean'
require 'time'

# Global variables
GO_MODULE="github.com/julb/go"
GO_APPS=[
    {:name => 'j3'}
]
GO_TARGET_ARCHS=[
    {:os => 'linux', :arch => 'amd64'},
    {:os => 'darwin', :arch => 'amd64'}
]
BUILD_DIR="build"

# Clean build directory
CLEAN.include "#{BUILD_DIR}/**/*"

# Format the source code
task :format do
    sh "go fmt ./...", verbose: false
end

# Lint the source code
task :lint do
    sh "golangci-lint run", verbose: false
end

# Run unit tests
namespace "test" do
    task :run do
        sh "go test ./...", verbose: false
    end

    task :coverage => ['clean'] do
        sh "
            mkdir -p #{BUILD_DIR}/test-coverage
            go test ./... -cover -coverprofile=#{BUILD_DIR}/test-coverage/report.out
            go tool cover -html=#{BUILD_DIR}/test-coverage/report.out -o #{BUILD_DIR}/test-coverage/report.html 
        "
    end
end

# Build binaries
desc 'build all projects'
task :build => ['clean'] do
    # Build vars.
    build_vars = {
        :time => Time.now.utc.iso8601,
        :version => "1.0.0",
        :git_revision => `git rev-parse HEAD`.strip!,
        :git_short_revision => `git rev-parse --short HEAD`.strip!
    }

    GO_APPS.each do |app|
        GO_TARGET_ARCHS.each do |target_arch|
            # log
            puts "build: #{app[:name]}:#{target_arch[:os]}:#{target_arch[:arch]}"

            # compile opts: output
            output_bin = "#{BUILD_DIR}/bin/#{target_arch[:os]}/#{target_arch[:arch]}/#{app[:name]}"

            # compile opts: ldflags.
            ld_flags="
                -s -w
                -X '#{GO_MODULE}/pkg/build.Group=me.julb.applications'
                -X '#{GO_MODULE}/pkg/build.Artifact=#{app[:name]}'
                -X '#{GO_MODULE}/pkg/build.Name=#{app[:name]}'
                -X '#{GO_MODULE}/pkg/build.Arch=#{target_arch[:os]}-#{target_arch[:arch]}'
                -X '#{GO_MODULE}/pkg/build.Version=#{build_vars[:version]}'
                -X '#{GO_MODULE}/pkg/build.Time=#{build_vars[:time]}'
                -X '#{GO_MODULE}/pkg/build.BuildVersion=#{build_vars[:version]}.#{build_vars[:git_short_revision]}'
                -X '#{GO_MODULE}/pkg/build.GitRevision=#{build_vars[:git_revision]}'
                -X '#{GO_MODULE}/pkg/build.GitShortRevision=#{build_vars[:git_short_revision]}'
            "
            
            # build: go binary
            sh "
                GOOS=#{target_arch[:os]}
                GOARCH=#{target_arch[:arch]}
                go build -ldflags=\"#{ld_flags}\" -o #{output_bin} cmd/#{app[:name]}/#{app[:name]}.go
            ", verbose: false

            # compute: md5
            open("#{output_bin}.md5", 'w') do |f|
                f.puts Digest::MD5.hexdigest File.read output_bin
            end

            # compute: sha256
            open("#{output_bin}.sha256", 'w') do |f|
                f.puts Digest::SHA256.hexdigest File.read output_bin
            end
        end
    end
end