require 'rake/clean'

# Clean task
CLEAN.include 'bin/**/*'

# Disable useless tasks.
#Rake::Task['release'].clear
#Rake::Task['install:local'].clear
#Rake::Task['install'].clear

desc 'build j3'
task :build_j3 do
    sh "go build -ldflags=\"-X 'github.com/julb/go/pkg/build.Version=1.0.0'\" -o bin/j3 cmd/j3/j3.go"
end

desc 'build everything'
task :build => ['clean', 'build_j3']
