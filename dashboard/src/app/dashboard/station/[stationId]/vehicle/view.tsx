"use client";

import { Stack } from "@mui/material";

import { useSearchParams } from "next/navigation";
import VehicleManageView from "./management/view";
import { MENU_TABS } from "../const";
import Transection from "./transaction";
import VehicleChart from "./chart/chart";
import { PropWithStationLocationId } from "../types";

export default function VehicleView({
  stationLocationId,
}: PropWithStationLocationId) {
  const menuParams = useSearchParams();
  const content = {
    [MENU_TABS[0]]: (
      <Stack spacing={1}>
        <VehicleChart />
        <Transection stationLocationId={stationLocationId} />
      </Stack>
    ),
    [MENU_TABS[1]]: <VehicleManageView />,
  };
  if (!content[menuParams.get("menu") as ""]) {
    return content[MENU_TABS[0]];
  }
  return content[menuParams.get("menu") as ""];
}
