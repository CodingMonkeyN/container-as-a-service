using k8s.Models;

namespace ContainerAsAService.Infrastructure.Utils;

public static class StatusUtil
{
    public static string MapStatus(V1ContainerState state)
    {
        if (state.Running != null)
        {
            return "Running";
        }

        if (state.Terminated != null)
        {
            return "Terminated";
        }

        if (state.Waiting != null)
        {
            return "Waiting";
        }

        return "Unknown";
    }
}
