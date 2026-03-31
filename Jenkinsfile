pipeline {
    agent any

    stages {
        stage('Install dependencies') {
            steps {
                bat 'python -m pip install --upgrade pip'
                bat 'python -m pip install selenium pytest webdriver-manager'
            }
        }

        stage('Run tests') {
            steps {
                bat 'python -m pytest tests --junitxml=report.xml'
            }
        }
    }

    post {
        always {
            junit 'report.xml'
        }
    }
}