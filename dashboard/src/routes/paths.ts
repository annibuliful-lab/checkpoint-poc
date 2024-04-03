export const paths = {
  dashboard: {
    notification: {
      root: "/dashboard/notification",
    },
    setting: {
      root: "/dashboard/setting",
    },
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
