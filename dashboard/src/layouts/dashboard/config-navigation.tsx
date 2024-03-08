import SvgColor from "@/components/svg-color";
import { paths } from "@/routes/paths";
import { useMemo } from "react";
import { SpaceDashboard } from "@mui/icons-material";
import { NavGroupProps } from "@/components/nav-section";
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
};

// ----------------------------------------------------------------------

export function useNavData() {
  const data: NavGroupProps[] = useMemo(
    () => [
      {
        items: [
          {
            title: "dashboard",
            path: paths.dashboard.root,
            icon: ICONS.dashboard,
          },
        ],
      },
    ],
    []
  );

  return data;
}
