import { Card } from "@mui/material";
import React from "react";
import ImeiTransectionTable from "./table/table";
import { PropWithStationLocationId } from "../types";

export default function Transection({
  stationLocationId,
}: PropWithStationLocationId) {
  return (
    <Card>
      <ImeiTransectionTable stationLocationId={stationLocationId} />
    </Card>
  );
}
