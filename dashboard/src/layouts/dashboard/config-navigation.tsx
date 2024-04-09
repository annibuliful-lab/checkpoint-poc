import SvgColor from "@/components/svg-color";
import { paths } from "@/routes/paths";
import { useMemo } from "react";
import {
  SpaceDashboard,
  Warehouse,
  PhoneIphone,
  NoCrash,
} from "@mui/icons-material";
import { NavGroupProps } from "@/components/nav-section";
import { useParams } from "next/navigation";
import Iconify from "@/components/iconify";
// ----------------------------------------------------------------------

const icon = (name: string) => (
  <SvgColor
    src={`/assets/icons/navbar/${name}.svg`}
    sx={{ width: 1, height: 1 }}
  />
  // OR
  // <Iconify icon="fluent:mail-24-filled" />
  // https://icon-sets.iconify.design/solar/
  // https://www.streamlinehq.com/icons
);

const ICONS = {
  dashboard: <SpaceDashboard />,
  station: <Warehouse />,
  imei: <Iconify icon={"material-symbols:p2p"} width={24} />,
  vehicle: <NoCrash />,
};

// ----------------------------------------------------------------------

export function useNavData() {
  const { stationId = "" } = useParams();
  const data: NavGroupProps[] = useMemo(() => {
    const items = [
      {
        path: paths.dashboard.station.root,
        icon: ICONS.station,
      },
    ];
    if (stationId) {
      items.push({
        path: paths.dashboard.station.join(`${stationId}`).dashboard.root,
        icon: ICONS.dashboard,
      });
      items.push({
        path: paths.dashboard.station.join(`${stationId}`).vehicle.root,
        icon: ICONS.vehicle,
      });
      items.push({
        path: paths.dashboard.station.join(`${stationId}`).imei.root,
        icon: ICONS.imei,
      });
    }
    return [
      {
        items,
      },
    ];
  }, [stationId]);

  return data;
}
