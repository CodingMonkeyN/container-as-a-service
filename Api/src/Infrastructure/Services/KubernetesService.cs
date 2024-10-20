using ContainerAsAService.Application.Common.Interfaces;
using ContainerAsAService.Application.Pods.Queries.GetPods;
using ContainerAsAService.Infrastructure.Utils;
using k8s;
using k8s.Models;

namespace ContainerAsAService.Infrastructure.Services;

public class KubernetesService(IKubernetes kubernetes) : IKubernetesService
{
    public async Task<Stream> GetPodLogsAsync(string podNamespace, string podName)
    {
        Stream? logResponse = await kubernetes.CoreV1.ReadNamespacedPodLogAsync(
            podName,
            podNamespace
        );
        return logResponse;
    }

    public async Task<IList<string>> GetClusterNamespacesAsync()
    {
        V1NamespaceList? namespaces = await kubernetes.CoreV1.ListNamespaceAsync();
        return namespaces.Items.Select(item => item.Metadata.Name).ToList();
    }

    public async Task<IList<PodDto>> GetPodsAsync()
    {
        V1PodList? pods = await kubernetes.CoreV1.ListPodForAllNamespacesAsync();
        if (pods is null)
        {
            return new List<PodDto>();
        }

        if (pods.Items is null)
        {
            return new List<PodDto>();
        }

        return pods.Items.Where(item => item.Spec.RuntimeClassName == "kata-qemu").Select(item => new PodDto
        {
            Name = item.Metadata.Name,
            Namespace = item.Metadata.NamespaceProperty,
            Cpu = item.Spec.Containers[0].Resources?.Limits?["cpu"].Value ?? "",
            Memory = item.Spec.Containers[0].Resources?.Limits?["memory"].Value ?? "",
            Port = item.Spec.Containers[0]?.Ports?[0].ContainerPort == null
                ? 0
                : item.Spec.Containers[0].Ports[0].ContainerPort,
            Ready = item.Status.ContainerStatuses[0].Ready,
            Status = StatusUtil.MapStatus(item.Status.ContainerStatuses[0].State)
        }).ToList();
    }
}
