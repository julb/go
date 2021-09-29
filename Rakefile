require 'digest'
require 'rake/clean'
require 'time'
require 'zip'

# Global variables
GO_MODULE="github.com/julb/go"
GO_APPS=[
#    {:name => 'j3', :type => 'cmd'},
#    {:name => 'prometheus-exporter', :type => 'cmd'},
#    {:name => 'hello-world', :type => 'aws/lambdas'},
    {:name => 'fun-gen-uuid', :type => 'aws/lambdas'},
    {:name => 'fun-gen-notification-content', :type => 'aws/lambdas'}
]
GO_APPS_TARGET_ARCHS=[
    {:os => 'linux', :arch => 'amd64'},
#    {:os => 'linux', :arch => 'arm64'},
#    {:os => 'darwin', :arch => 'amd64'}
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

    task :runv do
        sh "go test -v ./...", verbose: false
    end

    task :coverage => ['clean'] do
        sh "
            mkdir -p #{BUILD_DIR}/test-coverage
            go test ./... -cover -coverprofile=#{BUILD_DIR}/test-coverage/report.out
            go tool cover -html=#{BUILD_DIR}/test-coverage/report.out -o #{BUILD_DIR}/test-coverage/report.html 
        "
    end
end
task :test => ["test:run"]

# Run godoc serve
namespace "docs" do
    task :serve do
        sh "
            echo Open browser at: http://localhost:6060/pkg/#{GO_MODULE}
            godoc -http=:6060
        ", verbose: false
    end
end

# Build binaries
desc 'build all projects'
task :build, [:version] => ['clean'] do |t, args|
    args.with_defaults(:version => "__UNSET__")

    # Build vars.
    build_vars = {
        :time => Time.now.utc.iso8601,
        :group => "me.julb.applications",
        :version => args[:version],
        :git_revision => `git rev-parse HEAD`.strip!,
        :git_short_revision => `git rev-parse --short HEAD`.strip!
    }

    GO_APPS.each do |app|
        GO_APPS_TARGET_ARCHS.each do |target_arch|
            # log
            puts "build: #{app[:name]}:#{target_arch[:os]}:#{target_arch[:arch]}"

            # compile opts: output
            output_bin = "#{BUILD_DIR}/bin/#{target_arch[:os]}/#{target_arch[:arch]}/#{app[:name]}"

            # compile opts: ldflags.
            ldflags="
                -s -w
                -X '#{GO_MODULE}/pkg/build.Group=#{build_vars[:group]}'
                -X '#{GO_MODULE}/pkg/build.Artifact=#{app[:name]}'
                -X '#{GO_MODULE}/pkg/build.Name=#{app[:name]}'
                -X '#{GO_MODULE}/pkg/build.Arch=#{target_arch[:os]}-#{target_arch[:arch]}'
                -X '#{GO_MODULE}/pkg/build.Version=#{build_vars[:version]}'
                -X '#{GO_MODULE}/pkg/build.Time=#{build_vars[:time]}'
                -X '#{GO_MODULE}/pkg/build.BuildVersion=#{build_vars[:version]}.#{build_vars[:git_short_revision]}'
                -X '#{GO_MODULE}/pkg/build.GitRevision=#{build_vars[:git_revision]}'
                -X '#{GO_MODULE}/pkg/build.GitShortRevision=#{build_vars[:git_short_revision]}'
                -extldflags \"-static\"
            "
            
            # compile opts: tags.
            tags = "netgo"

            # build: go binary
            sh "
                export GOOS=#{target_arch[:os]}
                export GOARCH=#{target_arch[:arch]} 
                export CGO_ENABLED=0
                go build -tags #{tags} -ldflags=\"#{ldflags}\" -o #{output_bin} #{app[:type]}/#{app[:name]}/#{app[:name]}.go
            ", verbose: false

            # compute: md5
            open("#{output_bin}.md5", 'w') do |f|
                f.puts Digest::MD5.hexdigest File.read output_bin
            end

            # compute: sha256
            open("#{output_bin}.sha256", 'w') do |f|
                f.puts Digest::SHA256.hexdigest File.read output_bin
            end

            # generate zip file
            Zip::File.open("#{output_bin}.zip", create: true) do |zipfile|
                zipfile.add("main", output_bin)
                zipfile.add("#{app[:name]}", output_bin)
                zipfile.add("#{app[:name]}.md5", "#{output_bin}.md5")
                zipfile.add("#{app[:name]}.sha256", "#{output_bin}.sha256")
            end
        end
    end
end