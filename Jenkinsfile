pipeline {
  agent any
  stages {
    stage('Build') {
      parallel {
        stage('Build') {
          steps {
            sh 'echo "hello"'
          }
        }

        stage('Build 2') {
          steps {
            sh 'echo "HELLO 2"'
          }
        }

      }
    }

  }
}