package k8s

import (
	"context"

	"github.com/hirosassa/zerodriver"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"viesure.io/vs-initializer/pkg/config"
)

type (
	SecretsServiceInterface interface {
		Get(name string, opts metav1.GetOptions) (*v1.Secret, error)
		Update(secret *v1.Secret, opts metav1.UpdateOptions) (*v1.Secret, error)
	}

	SecretsService struct {
		namespace  string
		appContext context.Context
		clientset  *kubernetes.Clientset
		logger     *zerodriver.Logger
	}
)

func NewSecretsService(appContext context.Context) (*SecretsService, error) {

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &SecretsService{
		namespace:  viper.GetString(config.AppNamespaceVarName),
		appContext: appContext,
		clientset:  clientset,
		logger:     viper.Get(config.LoggerVarName).(*zerodriver.Logger),
	}, nil

}

func (svc *SecretsService) Get(name string, options metav1.GetOptions) (*v1.Secret, error) {

	result, err := svc.clientset.
		CoreV1().
		Secrets(svc.namespace).
		Get(svc.appContext, name, options)

	if err != nil {
		return nil, err
	}

	return result, nil

}

func (svc *SecretsService) Update(secret *v1.Secret, options metav1.UpdateOptions) (*v1.Secret, error) {

	result, err := svc.clientset.
		CoreV1().
		Secrets(svc.namespace).
		Update(svc.appContext, secret, options)

	if err != nil {
		return nil, err
	}

	return result, nil

}
