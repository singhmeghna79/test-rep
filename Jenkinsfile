def ORG = "mayadataio"
def REPO = "graph-reporter"
def DOCKER_HUB_REPO = "https://index.docker.io/v1/"
def DOCKER_IMAGE = ""
def STAGING_KEYPATH="~jenkins/.ssh/id_rsa_control_node_staging"
def PROD_KEYPATH="~jenkins/.ssh/id_rsa_control_node_production"
def PREPROD_KEYPATH="~jenkins/.ssh/id_rsa_control_node_preproduction"
def CONTROL_NODE="35.225.61.42"
env.user="atulabhi"
env.pass="ka879707"
pipeline {
    agent any 
    stages {
        stage('Build Image') {
            steps {
                script {
                    def root = tool name: '1.9.2', type: 'go'
                    
                    /**
                     * Add in GOPATH, GOROOT, GOBIN to the environment and add go to the path for jenkins.
                     * The environment variables are needed for godep to properly work and adding go to the path is required to
                     */
                    withEnv(["GOPATH=${env.WORKSPACE}/go", "GOROOT=${root}", "GOBIN=${root}/bin", "PATH+GO=${root}/bin"]) {
                        sh "mkdir -p ${env.WORKSPACE}/go/src/github.com/mayadata-io/graph-reporter"
                        echo "Building Go Code"
                        sh """
                        rsync -r --exclude 'go' ${WORKSPACE}/* ${GOPATH}/src/github.com/mayadata-io/graph-reporter
                        cd ${env.WORKSPACE}/go/src/github.com/mayadata-io/graph-reporter
                        echo ${env.BRANCH_NAME} > BRANCH.txt
                        go get -u github.com/golang/dep/cmd/dep
                        dep ensure
                        make build 
                        mv tag.txt ${env.WORKSPACE}
                        """

                    }
                }
            }
        }
        stage('Push Image') {
            steps {
                script {
		             docker.withRegistry('https://registry.hub.docker.com', 'ddc3fdf7-5611-4d47-a8ab-d0ea7624671a') {
                            if (env.BRANCH_NAME == 'staging'  ||env.BRANCH_NAME.startsWith('alpha-r')) {
		                       echo "Pushing the image with the tag..."
                               tag=sh(
		                                returnStdout: true,
				                        script: "cat tag.txt"
				                      ).trim()
                                        
                               sh "docker build -t mayadataio/graph-reporter:${env.BRANCH_NAME}-${tag} ."        
                               sh "docker push mayadataio/graph-reporter:${env.BRANCH_NAME}-${tag}"
                               
				            
                            } else if (env.BRANCH_NAME == 'master')  {
                               echo "Pushing the image with the tag..."
                               tag=sh(
		                                returnStdout: true,
				                        script: "cat tag.txt"
				                      ).trim() 
                               sh "docker build -t mayadataio/graph-reporter:${tag} ."       
			                   sh "docker push mayadataio/graph-reporter:${tag}"				            			            
                            } else {
			                   echo "WARNING: Not pushing Image"
                        }       
                    }
                }
            }
        }
        stage('Deploy on the related k8s cluster') {
           steps {
               script {
                   /**
                     * After Image was successfuly built and push to docker here it will be deployed to respective k8s cluster.
                     */
                   if (env.BRANCH_NAME == 'staging') {
                       // Deploy to staging cluster
                       echo "${env.BRANCH_NAME}-${tag}"
                       sh "ssh -i ${STAGING_KEYPATH} staging@${CONTROL_NODE} \" /home/staging/install.sh graph-reporter \"${env.BRANCH_NAME}-${tag}\"\""
                   } else if (env.BRANCH_NAME == 'master') {
                       // Deploy to production cluster
                       sh "ssh -i ${PROD_KEYPATH} production@${CONTROL_NODE} \" /home/production/install.sh graph-reporter \"${tag}\"\""
                   } else if(env.BRANCH_NAME.startsWith('alpha-r') || env.BRANCH_NAME == 'release') {
                       // Deploy to pre-production cluster
                       sh "ssh -i ${PREPROD_KEYPATH} preproduction@${CONTROL_NODE} \" /home/preproduction/install.sh graph-reporter \"${env.BRANCH_NAME}-${tag}\"\""
                   } else {
                       echo "Not sure what to do with this branch. So not deploying. May be dev branch ?"
                   }
                }
             }
         }
        stage ('Adding git-tag for master') {
          steps {
            script {
                /**
                     * Add in GOPATH, GOROOT, GOBIN to the environment and add go to the path for jenkins.
                     * The environment variables are needed for godep to properly work and adding go to the path is required to
                     */
               if (env.BRANCH_NAME == 'master') {
                 sh """
                  git tag -fa "${tag}" -m "Release of ${tag}"
                  """
                  sh "git tag -l"
                  sh """
                    git push https://${env.user}:${env.pass}@github.com/mayadata-io/graph-reporter.git --tag
                     """
             }
          }
        }
     } 
        
   } 
    /**
     * This post step will always execute regardless of a build failing or passing to clean up the setting that allows godep
     * to install private dependencies from Github. When using `checkout scm` or the default Jenkins clone step for a git
     * multibranch pipeline this undo change is needed. If the url change is not undone it will fail subsequent builds because the
     * Jenkins Git plugin will fail to clone the repository correctly.
     */
    post {
        always {
            echo 'This will always run'
            deleteDir()
        }
        success {
            echo 'This will run only if successful'
            slackSend channel: '#jenkins-builds',
                   color: 'good',
                   message: "The pipeline ${currentBuild.fullDisplayName} completed successfully :dance: :thumbsup: "
            
	}
        failure {
            echo 'This will run only if failed'
            slackSend channel: '#jenkins-builds',
                  color: 'RED',
                  message: "The pipeline ${currentBuild.fullDisplayName} failed. :scream_cat: :japanese_goblin: "

        }
        unstable {
            echo 'This will run only if the run was marked as unstable'
            slackSend channel: '#jenkins-builds',
                   color: 'good',
                   message: "The pipeline ${currentBuild.fullDisplayName} is unstable :scream_cat: :japanese_goblin: "

	}
        changed {
/*            slackSend channel: '#jenkins-builds',
                   color: 'good',
                   message: "Build ${currentBuild.fullDisplayName} is now stable :dance: :thumbsup: "
            echo 'This will run only if the state of the Pipeline has changed'
*/            echo 'For example, if the Pipeline was previously failing but is now successful'
        }
    }
}