export const paths = {
  dashboard: {
    root: "/dashboard",
    station: {
      root: "/dashboard/station",
      join: (stationId: string) => ({
        imei: {
          root: `/dashboard/station/${stationId}/imei`,
        },
        vehicle: {
          root: `/dashboard/station/${stationId}/vehicle`,
        },
        dashboard: {
          root: `/dashboard/station/${stationId}/dashboard`,
        },
      }),
    },
  },
  auth: {
    login: "/auth/login",
  },
};
