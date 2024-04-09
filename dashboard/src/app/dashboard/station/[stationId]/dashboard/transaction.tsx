import { Card } from "@mui/material";
import React from "react";
import VehicleTransectionTable from "./table/table";
import { PropWithStationLocationId } from "../types";

export default function Transection({
  stationLocationId,
}: PropWithStationLocationId) {
  return (
    <Card>
      <VehicleTransectionTable stationLocationId={stationLocationId} />
    </Card>
  );
}
