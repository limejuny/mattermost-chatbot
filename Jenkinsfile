pipeline {
    agent {
        kubernetes {
            inheritFrom 'go1.18'
        }
    }

    stages {
        stage("Build") {
            container('go') {
                script {
                    sh """
                        make build
                    """
                }
            }
        }

        stage("Check Binary") {
            steps {
                script {
                    sh """
                        ls -lh bin/
                    """
                }
            }
        }
    }
}
