"use client";

import { AuthGuard } from "@/auth/guard";
import DashboardLayout from "@/layouts/dashboard";
import { Add, CameraOutdoor, Close } from "@mui/icons-material";
import { Box, Divider, IconButton, Stack, Typography } from "@mui/material";
import { useAtom } from "jotai";
import { stationHistoriesAtom } from "../store";
import { useCallback, useMemo } from "react";
import { usePathname, useRouter, useSearchParams } from "@/routes/hooks";
import { paths } from "@/routes/paths";
import { MENU_TABS } from "./const";
// ----------------------------------------------------------------------

type Props = {
  children: React.ReactNode;
  params: { stationId: string };
};

export default function Layout({ children, params }: Props) {
  const menuParams = useSearchParams();
  const [stationHistories, setStationHistories] = useAtom(stationHistoriesAtom);
  const pathname = usePathname();
  const isDashboardScreen = useMemo(
    () =>
      paths.dashboard.station.join(params.stationId).dashboard.root ===
      pathname,
    [params.stationId, pathname]
  );

  const route = useRouter();
  const onChangeStation = useCallback(
    (id: string) => {
      route.push(paths.dashboard.station.join(id).dashboard.root);
    },
    [route]
  );

  const onAddTags = useCallback(() => {
    route.push(paths.dashboard.station.root);
  }, [route]);

  const onChangeMenu = useCallback(
    (menu: string) => {
      route.push(`?menu=${menu}`);
    },
    [route]
  );

  const selectedMenu = useCallback(() => {
    if (menuParams.get("menu")) {
      return menuParams.get("menu");
    }
    return MENU_TABS[0];
  }, [menuParams]);
  const stationHistoryIds = useMemo(
    () => Object.keys(stationHistories),
    [stationHistories]
  );
  const onCloseStation = useCallback(
    (index: number) => {
      let initId = stationHistoryIds[index];
      setStationHistories((prevs) => {
        delete prevs[initId];
        return { ...prevs };
      });
      if (initId === params.stationId) {
        let redirectId = initId;
        const endStationTabs = stationHistoryIds.length === index + 1;
        if (index > 0) {
          if (endStationTabs) {
            redirectId = stationHistoryIds[index - 1];
          } else {
            redirectId = stationHistoryIds[index + 1];
          }
          route.push(paths.dashboard.station.join(redirectId).dashboard.root);
        } else {
          route.push(paths.dashboard.station.root);
        }
      }
    },
    [params.stationId, route, setStationHistories, stationHistoryIds]
  );
  const stationHistoryTabs = (
    <Stack
      direction={"row"}
      alignItems={"center"}
      divider={<Divider flexItem orientation="vertical" />}
      bgcolor={"#FFFFFF"}
    >
      <Stack direction={"row"} spacing={1} px={3} py={1}>
        <CameraOutdoor />
        <Typography>Station site</Typography>
      </Stack>
      {stationHistoryIds.map((id, index) => (
        <Stack
          direction={"row"}
          spacing={1}
          p={1}
          px={2}
          key={id}
          component={Box}
          justifyContent={"space-between"}
          alignItems={"center"}
          flex={stationHistoryIds.length > 6 ? 1 : undefined}
          sx={{
            height: 40,
            minWidth: 140,
            background: id === params.stationId ? "#EEF4FF" : undefined,
            borderBottom:
              id === params.stationId
                ? "2px solid #142FE1"
                : "2px solid transparent",
          }}
        >
          <Typography
            variant="subtitle2"
            fontSize={12}
            flex={1}
            onClick={() => onChangeStation(id)}
            sx={{
              cursor: "pointer",
            }}
          >
            {stationHistories[id]?.title}
          </Typography>
          <IconButton size="small" onClick={() => onCloseStation(index)}>
            <Close sx={{ fontSize: 12 }} />
          </IconButton>
        </Stack>
      ))}
      <Box px={1}>
        <IconButton size="small" onClick={onAddTags}>
          <Add fontSize="small" />
        </IconButton>
      </Box>
    </Stack>
  );
  const menuTabs = (
    <Stack direction={"row"} alignItems={"center"} bgcolor={"#FFFFFF"}>
      {MENU_TABS.map((menu) => (
        <Stack
          direction={"row"}
          spacing={1}
          px={3}
          py={0.5}
          key={menu}
          component={Box}
          onClick={() => onChangeMenu(menu)}
          sx={{
            cursor: "pointer",
            borderBottom:
              selectedMenu() === menu
                ? "2px solid #142FE1"
                : "2px solid transparent",
          }}
        >
          <Typography>{menu}</Typography>
        </Stack>
      ))}
    </Stack>
  );
  return (
    <Stack bgcolor={"#F9FAFA"} minHeight={"100vh"}>
      {stationHistoryTabs}
      <Divider />
      {!isDashboardScreen && (
        <>
          {menuTabs}
          <Divider />
        </>
      )}
      <Box p={isDashboardScreen ? 0.5 : 2}>{children}</Box>
    </Stack>
  );
}
