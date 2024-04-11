import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import MenuItem from "@mui/material/MenuItem";
import TableRow from "@mui/material/TableRow";

import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";

import { useBoolean } from "@/hooks/use-boolean";
import Iconify from "@/components/iconify";

import CustomPopover, { usePopover } from "@/components/custom-popover";
import Label from "@/components/label";
import { ConfirmDialog } from "@/components/custom-dialog";
import { Aod, ChevronRight, Edit, SimCard } from "@mui/icons-material";
import {
  Alert,
  Badge,
  Box,
  Tooltip,
  TooltipProps,
  Typography,
  styled,
  tooltipClasses,
} from "@mui/material";
import { useRouter } from "@/routes/hooks";
import { VehicleTransection } from "./types";
import { fDate } from "@/utils/format-time";
import { DevicePermittedLabel } from "@/apollo-client";
import Image from "@/components/image";
import { RouterLink } from "@/routes/components";
import { paths } from "@/routes/paths";

// ----------------------------------------------------------------------

type Props = {
  row: VehicleTransection;
  onDeleteRow: VoidFunction;
  onEdit: VoidFunction;
};

const HtmlTooltip = styled(({ className, ...props }: TooltipProps) => (
  <Tooltip {...props} classes={{ popper: className }} />
))(({ theme }) => ({
  [`& .${tooltipClasses.tooltip}`]: {
    backgroundColor: "#f5f5f9",
    color: "rgba(0, 0, 0, 0.87)",
    maxWidth: 220,
    fontSize: theme.typography.pxToRem(12),
    border: "1px solid #dadde9",
  },
}));

export default function VehicleTableRow({ row, onDeleteRow, onEdit }: Props) {
  const {
    id,
    remark,
    tags,
    arrivalTime,
    licensePlate,
    brand,
    stationSite,
    imei,
    imsi,
    color,
    vehicle,
  } = row;
  const confirm = useBoolean();
  const popover = usePopover();
  const route = useRouter();
  const isBlackList = licensePlate.status === DevicePermittedLabel.Blacklist;
  const isWhiteList = licensePlate.status === DevicePermittedLabel.Whitelist;
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
          <Tooltip
            arrow
            placement="bottom-end"
            title={
              <Box sx={{ width: 120 }}>
                <Image
                  width={120}
                  src={licensePlate.image ?? "/assets/images/license.jpg"}
                  alt={licensePlate.license}
                />
              </Box>
            }
          >
            <Typography
              fontSize={14}
              sx={{ cursor: "default" }}
              color={isBlackList || isWhiteList ? "white" : undefined}
            >
              {licensePlate.license}
            </Typography>
          </Tooltip>
        </TableCell>
        <TableCell>
          <Typography
            fontSize={14}
            color={isBlackList || isWhiteList ? "white" : undefined}
          >
            {licensePlate.type}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            fontSize={14}
            color={isBlackList || isWhiteList ? "white" : undefined}
          >
            {brand}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            fontSize={14}
            color={isBlackList || isWhiteList ? "white" : undefined}
          >
            {vehicle.type}
          </Typography>
        </TableCell>
        <TableCell>
          <Stack direction={"row"} spacing={1} alignItems={"center"}>
            <Box
              height={15}
              width={15}
              borderRadius={"50%"}
              bgcolor={color.code}
            />
            <Typography
              fontSize={14}
              color={isBlackList || isWhiteList ? "white" : undefined}
            >
              {color.name}
            </Typography>
          </Stack>
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
          <Tooltip
            arrow
            placement="bottom-end"
            title={
              <Stack spacing={1} sx={{ width: 180 }} pt={1}>
                {imei.list.map((ime) => (
                  <Stack key={ime} direction={"row"} spacing={0.5}>
                    <Aod fontSize="small" />
                    <Typography fontSize={14}>{ime}</Typography>
                  </Stack>
                ))}
                <Button
                  size="small"
                  LinkComponent={RouterLink}
                  href={"/"}
                  endIcon={<ChevronRight sx={{ fontSize: 14 }} />}
                  variant="text"
                  color="info"
                >
                  <Typography color={"info.light"} fontSize={10} align="center">
                    Go to IMEI/IMSI Dashboard
                  </Typography>
                </Button>
              </Stack>
            }
          >
            <Typography
              fontSize={14}
              sx={{ cursor: "default" }}
              color={isBlackList || isWhiteList ? "white" : undefined}
              align="center"
            >
              {imei.total}
            </Typography>
          </Tooltip>
        </TableCell>
        <TableCell>
          <Tooltip
            arrow
            placement="bottom-end"
            title={
              <Stack spacing={1} sx={{ width: 180 }} pt={1}>
                {imsi.list.map((ime) => (
                  <Stack key={ime} direction={"row"} spacing={0.5}>
                    <SimCard fontSize="small" />
                    <Typography fontSize={14}>{ime}</Typography>
                  </Stack>
                ))}

                <Button
                  size="small"
                  LinkComponent={RouterLink}
                  href={"/"}
                  endIcon={<ChevronRight sx={{ fontSize: 14 }} />}
                  variant="text"
                  color="info"
                >
                  <Typography color={"info.light"} fontSize={10} align="center">
                    Go to IMEI/IMSI Dashboard
                  </Typography>
                </Button>
              </Stack>
            }
          >
            <Typography
              fontSize={14}
              sx={{ cursor: "default" }}
              color={isBlackList || isWhiteList ? "white" : undefined}
              align="center"
            >
              {imsi.total}
            </Typography>
          </Tooltip>
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
        <TableCell>
          <IconButton color="warning">
            <Iconify icon={"tabler:alert-circle"} />
          </IconButton>
        </TableCell>
      </TableRow>
    </>
  );
}
