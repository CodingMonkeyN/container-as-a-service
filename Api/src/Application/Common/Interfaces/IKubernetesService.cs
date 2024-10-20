using ContainerAsAService.Application.Pods.Queries.GetPods;

namespace ContainerAsAService.Application.Common.Interfaces;

public interface IKubernetesService
{
    Task<Stream> GetPodLogsAsync(string podNamespace, string podName);

    Task<IList<string>> GetClusterNamespacesAsync();

    Task<IList<PodDto>> GetPodsAsync();
}
