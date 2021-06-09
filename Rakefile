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

# Clean task
CLEAN.include 'bin/**/*'

# Disable useless tasks.
#Rake::Task['release'].clear
#Rake::Task['install:local'].clear
#Rake::Task['install'].clear

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
            output_bin = "bin/#{target_arch[:os]}-#{target_arch[:arch]}/#{app[:name]}"

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
            open("#{output_bin}.md5sum", 'w') do |f|
                f.puts Digest::MD5.hexdigest File.read output_bin
            end

            # compute: sha256
            open("#{output_bin}.sha256sum", 'w') do |f|
                f.puts Digest::SHA256.hexdigest File.read output_bin
            end
        end
    end
end

desc 'format source code'
task :format do
    sh "find . -name '*.go' -exec go fmt {} \\;", verbose: false
end

desc 'lint source code'
task :lint do
    sh "golangci-lint run"
end