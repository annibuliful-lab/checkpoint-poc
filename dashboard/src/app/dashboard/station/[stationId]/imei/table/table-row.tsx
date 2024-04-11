import Stack from "@mui/material/Stack";
import TableRow from "@mui/material/TableRow";
import TableCell from "@mui/material/TableCell";
import { useBoolean } from "@/hooks/use-boolean";
import Label from "@/components/label";
import { Typography } from "@mui/material";
import { ImeiImsiTransaction } from "./types";
import { fDate } from "@/utils/format-time";
import { DevicePermittedLabel } from "@/apollo-client";

// ----------------------------------------------------------------------

type Props = {
  row: ImeiImsiTransaction;
  onDeleteRow: VoidFunction;
  onEdit: VoidFunction;
};

export default function ImeiTableRow({ row, onDeleteRow, onEdit }: Props) {
  const {
    id,
    tags,
    arrivalTime,
    imei,
    imsi,
    phoneModel,
    licensePlate,
    stationSite,
    status,
  } = row;

  const isBlackList = status === DevicePermittedLabel.Blacklist;
  const isWhiteList = status === DevicePermittedLabel.Whitelist;
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
            fontSize={14}
            color={isBlackList || isWhiteList ? "white" : undefined}
          >
            {fDate(arrivalTime)}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            fontSize={14}
            color={isBlackList || isWhiteList ? "white" : undefined}
          >
            {imei}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            fontSize={14}
            color={isBlackList || isWhiteList ? "white" : undefined}
          >
            {imsi}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            fontSize={14}
            color={isBlackList || isWhiteList ? "white" : undefined}
          >
            {phoneModel}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            fontSize={14}
            color={isBlackList || isWhiteList ? "white" : undefined}
          >
            {licensePlate}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            fontSize={14}
            color={isBlackList || isWhiteList ? "white" : undefined}
          >
            {stationSite}
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
