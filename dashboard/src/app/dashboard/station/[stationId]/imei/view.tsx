"use client";

import { Stack } from "@mui/material";

import { useSearchParams } from "next/navigation";
import ImeiManageView from "./management/view";
import { MENU_TABS } from "../const";
import Transection from "./transaction";
import ImeiChart from "./chart/chart";
import { PropWithStationLocationId } from "../types";

export default function ImeiView({
  stationLocationId,
}: PropWithStationLocationId) {
  const menuParams = useSearchParams();
  const content = {
    [MENU_TABS[0]]: (
      <Stack spacing={1}>
        <ImeiChart />
        <Transection stationLocationId={stationLocationId} />
      </Stack>
    ),
    [MENU_TABS[1]]: <ImeiManageView stationLocationId={stationLocationId} />,
  };
  if (!content[menuParams.get("menu") as ""]) {
    return content[MENU_TABS[0]];
  }
  return content[menuParams.get("menu") as ""];
}
