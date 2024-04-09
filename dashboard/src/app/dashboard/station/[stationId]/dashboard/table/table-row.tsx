import Stack from "@mui/material/Stack";
import TableRow from "@mui/material/TableRow";

import TableCell from "@mui/material/TableCell";

import Label from "@/components/label";
import { Typography } from "@mui/material";
import { fDateTime } from "@/utils/format-time";
import { DevicePermittedLabel } from "@/apollo-client";
import { StationDashboardTransaction } from "./types";

// ----------------------------------------------------------------------

type Props = {
  row: StationDashboardTransaction;
  onDeleteRow: VoidFunction;
  onEdit: VoidFunction;
};

export default function VehicleTableRow({ row, onDeleteRow, onEdit }: Props) {
  const { id, tags, arrivalTime, imei, imsi, phoneModel } = row;

  const isBlackList = [imei?.permittedLabel, imsi?.permittedLabel].includes(
    DevicePermittedLabel.Blacklist
  );
  const isWhiteList = [imei?.permittedLabel, imsi?.permittedLabel].includes(
    DevicePermittedLabel.Whitelist
  );
  return (
    <>
      <TableRow
        sx={{
          bgcolor:
            (isBlackList && "error.dark") ||
            (isWhiteList && "success.dark") ||
            undefined,
        }}
      >
        <TableCell>
          <Typography
            color={isBlackList || isWhiteList ? "white" : undefined}
            fontSize={14}
          >
            {fDateTime(arrivalTime)}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            color={isBlackList || isWhiteList ? "white" : undefined}
            fontSize={14}
          >
            {imei?.imei}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            color={isBlackList || isWhiteList ? "white" : undefined}
            fontSize={14}
          >
            {imsi?.imsi}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            color={isBlackList || isWhiteList ? "white" : undefined}
            fontSize={14}
          >
            {phoneModel}
          </Typography>
        </TableCell>
        <TableCell>
          <Stack
            alignItems={"center"}
            spacing={1}
            direction={"row"}
            flexWrap={"wrap"}
          >
            {tags?.map((tag) => (
              <Label
                key={tag.tag}
                variant={isBlackList || isWhiteList ? "filled" : "soft"}
                color="warning"
              >
                {tag.tag}
              </Label>
            ))}
          </Stack>
        </TableCell>
      </TableRow>
    </>
  );
}
