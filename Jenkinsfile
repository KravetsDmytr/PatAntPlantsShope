pipeline {
    agent any

    stages {
        stage('Setup Python venv') {
            steps {
                bat '''
@echo off
where python >nul 2>nul
if %ERRORLEVEL%==0 (
  set PYTHON_EXE=python
) else (
  set PYTHON_EXE=py -3
)
%PYTHON_EXE% -m venv .venv
'''
                bat '.venv\\Scripts\\python -m pip install --upgrade pip'
                bat '.venv\\Scripts\\python -m pip install -r requirements.txt'
            }
        }

        stage('Start API for tests') {
            steps {
                bat 'docker compose up --build -d postgres web'
            }
        }

        stage('Run tests') {
            steps {
                bat '.venv\\Scripts\\python -m pytest tests --junitxml=report.xml -v'
            }
        }
    }

    post {
        always {
            junit allowEmptyResults: true, testResults: 'report.xml'
            bat 'docker compose down'
        }
    }
}