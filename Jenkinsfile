pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh 'echo "hello"'
      }
    }

    stage('Next') {
      parallel {
        stage('PWD') {
          steps {
            sh 'pwd'
          }
        }

        stage('LS') {
          steps {
            sh 'ls'
          }
        }

        stage('Message') {
          steps {
            echo 'message'
          }
        }

      }
    }

  }
}