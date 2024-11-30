import { QueryClient } from "@tanstack/react-query"

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            refetchIntervalInBackground: true,
            refetchOnWindowFocus: 'always',
            refetchOnMount: true,
            refetchInterval: 1000 * 15,
        },
    },
})

export default queryClient

