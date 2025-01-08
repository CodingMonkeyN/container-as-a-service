using k8s.Models;

namespace Api.Utils;

public static class Utils
{
    public static string MapStatus(V1ContainerState state)
    {
        if (state.Running != null)
            return "Running";
        if (state.Terminated != null)
            return "Terminated";
        if (state.Waiting != null)
            return "Waiting";
        return "Unknown";
    }
}