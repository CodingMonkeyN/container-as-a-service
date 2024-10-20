using ContainerAsAService.Application.Common.Interfaces;
using ContainerAsAService.Infrastructure.Services;
using k8s;
using Microsoft.Extensions.Hosting;

namespace Microsoft.Extensions.DependencyInjection;

public static class DependencyInjection
{
    public static void AddInfrastructureServices(this IHostApplicationBuilder builder)
    {
        KubernetesClientConfiguration? config = KubernetesClientConfiguration.BuildConfigFromConfigFile();
        Kubernetes client = new Kubernetes(config);
        builder.Services.AddSingleton(TimeProvider.System);
        builder.Services.AddScoped<IKubernetesService, KubernetesService>();
        builder.Services.AddScoped<IContainerDeploymentService, ContainerDeploymentService>();
        builder.Services.AddSingleton<IKubernetes>(builder => client);
    }
}
