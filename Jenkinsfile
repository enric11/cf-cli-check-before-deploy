pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh '''echo "hello"
. ~/.nvm/nvm.sh
npm -v
mbt -v'''
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