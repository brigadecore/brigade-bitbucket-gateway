if 'ENABLE_NGROK_EXTENSION' in os.environ and os.environ['ENABLE_NGROK_EXTENSION'] == '1':
  v1alpha1.extension_repo(
    name = 'default',
    url = 'https://github.com/tilt-dev/tilt-extensions'
  )
  v1alpha1.extension(name = 'ngrok', repo_name = 'default', repo_path = 'ngrok')

load('ext://min_k8s_version', 'min_k8s_version')
min_k8s_version('1.18.0')

trigger_mode(TRIGGER_MODE_MANUAL)

load('ext://namespace', 'namespace_create')
namespace_create('brigade-bitbucket-gateway')
k8s_resource(
  new_name = 'namespace',
  objects = ['brigade-bitbucket-gateway:namespace'],
  labels = ['brigade-bitbucket-gateway']
)

docker_build(
  'brigadecore/brigade-bitbucket-gateway', '.',
  only = [
    'internal/',
    'config.go',
    'go.mod',
    'go.sum',
    'main.go'
  ],
  ignore = ['**/*_test.go']
)
k8s_resource(
  workload = 'brigade-bitbucket-gateway',
  new_name = 'gateway',
  port_forwards = '31700:8080',
  labels = ['brigade-bitbucket-gateway']
)
k8s_resource(
  workload = 'gateway',
  objects = ['brigade-bitbucket-gateway:secret']
)

k8s_yaml(
  helm(
    './charts/brigade-bitbucket-gateway',
    name = 'brigade-bitbucket-gateway',
    namespace = 'brigade-bitbucket-gateway',
    set = [
      'brigade.apiToken=' + os.environ['BRIGADE_API_TOKEN'],
      'service.type=NodePort',
      'tls.enabled=false'
    ]
  )
)
