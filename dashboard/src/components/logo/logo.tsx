import { forwardRef } from "react";

import Link from "@mui/material/Link";
import { useTheme } from "@mui/material/styles";
import Box, { BoxProps } from "@mui/material/Box";

import { RouterLink } from "@/routes/components";
import { paths } from "@/routes/paths";

// ----------------------------------------------------------------------

export interface LogoProps extends BoxProps {
  disabledLink?: boolean;
}

// eslint-disable-next-line react/display-name
const Logo = forwardRef<HTMLDivElement, LogoProps>(
  ({ disabledLink = false, sx, ...other }, ref) => {
    const theme = useTheme();

    const PRIMARY_LIGHT = theme.palette.primary.light;

    const PRIMARY_MAIN = theme.palette.primary.main;

    const PRIMARY_DARK = theme.palette.primary.dark;

    // OR using local (public folder)
    // -------------------------------------------------------
    const logo = (
      <Box
        component="img"
        src={"/assets/logos/logo.svg"}
        sx={{ width: 40, height: 40, cursor: "pointer", ...sx }}
      />
    );

    // const logo = (
    //   <Box
    //     ref={ref}
    //     component="div"
    //     sx={{
    //       width: 40,
    //       height: 40,
    //       display: "inline-flex",
    //       ...sx,
    //     }}
    //     {...other}
    //   >
    //     <svg
    //       width="42"
    //       height="32"
    //       viewBox="0 0 42 32"
    //       fill="none"
    //       xmlns="http://www.w3.org/2000/svg"
    //     >
    //       <g filter="url(#filter0_d_139_452)">
    //         <path
    //           d="M3 3H15.3819V5.32161H8.62814V22.6784H15.3819V25H3V3Z"
    //           fill="url(#paint0_linear_139_452)"
    //         />
    //         <path
    //           d="M25.6197 19L35 15.9045V24.871H25.6197V19Z"
    //           fill="url(#paint1_linear_139_452)"
    //         />
    //         <path
    //           fill-rule="evenodd"
    //           clip-rule="evenodd"
    //           d="M17.3116 3V25H22.9396V16.4874L35 12.397V4.65829L31.3015 3H17.3116ZM31.3015 10.8738L22.9396 13.6376V5.34618L31.3015 5.32161V10.8738Z"
    //           fill="url(#paint2_linear_139_452)"
    //         />
    //         <path
    //           d="M3 3H15.3819V5.32161H8.62814V22.6784H15.3819V25H3V3Z"
    //           stroke="url(#paint3_linear_139_452)"
    //         />
    //         <path
    //           d="M25.6197 19L35 15.9045V24.871H25.6197V19Z"
    //           stroke="url(#paint4_linear_139_452)"
    //         />
    //         <path
    //           fill-rule="evenodd"
    //           clip-rule="evenodd"
    //           d="M17.3116 3V25H22.9396V16.4874L35 12.397V4.65829L31.3015 3H17.3116ZM31.3015 10.8738L22.9396 13.6376V5.34618L31.3015 5.32161V10.8738Z"
    //           stroke="url(#paint5_linear_139_452)"
    //         />
    //       </g>
    //       <defs>
    //         <filter
    //           id="filter0_d_139_452"
    //           x="0.5"
    //           y="0.5"
    //           width="41"
    //           height="31"
    //           filterUnits="userSpaceOnUse"
    //           color-interpolation-filters="sRGB"
    //         >
    //           <feFlood flood-opacity="0" result="BackgroundImageFix" />
    //           <feColorMatrix
    //             in="SourceAlpha"
    //             type="matrix"
    //             values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0"
    //             result="hardAlpha"
    //           />
    //           <feOffset dx="2" dy="2" />
    //           <feGaussianBlur stdDeviation="2" />
    //           <feComposite in2="hardAlpha" operator="out" />
    //           <feColorMatrix
    //             type="matrix"
    //             values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.1 0"
    //           />
    //           <feBlend
    //             mode="normal"
    //             in2="BackgroundImageFix"
    //             result="effect1_dropShadow_139_452"
    //           />
    //           <feBlend
    //             mode="normal"
    //             in="SourceGraphic"
    //             in2="effect1_dropShadow_139_452"
    //             result="shape"
    //           />
    //         </filter>
    //         <linearGradient
    //           id="paint0_linear_139_452"
    //           x1="19"
    //           y1="3"
    //           x2="19"
    //           y2="25"
    //           gradientUnits="userSpaceOnUse"
    //         >
    //           <stop stop-color="#254EDB" />
    //           <stop offset="1" stop-color="#142A75" />
    //         </linearGradient>
    //         <linearGradient
    //           id="paint1_linear_139_452"
    //           x1="19"
    //           y1="3"
    //           x2="19"
    //           y2="25"
    //           gradientUnits="userSpaceOnUse"
    //         >
    //           <stop stop-color="#254EDB" />
    //           <stop offset="1" stop-color="#142A75" />
    //         </linearGradient>
    //         <linearGradient
    //           id="paint2_linear_139_452"
    //           x1="19"
    //           y1="3"
    //           x2="19"
    //           y2="25"
    //           gradientUnits="userSpaceOnUse"
    //         >
    //           <stop stop-color="#254EDB" />
    //           <stop offset="1" stop-color="#142A75" />
    //         </linearGradient>
    //         <linearGradient
    //           id="paint3_linear_139_452"
    //           x1="19"
    //           y1="3"
    //           x2="19"
    //           y2="25"
    //           gradientUnits="userSpaceOnUse"
    //         >
    //           <stop stop-color="#254EDB" />
    //           <stop offset="1" stop-color="#142A75" />
    //         </linearGradient>
    //         <linearGradient
    //           id="paint4_linear_139_452"
    //           x1="19"
    //           y1="3"
    //           x2="19"
    //           y2="25"
    //           gradientUnits="userSpaceOnUse"
    //         >
    //           <stop stop-color="#254EDB" />
    //           <stop offset="1" stop-color="#142A75" />
    //         </linearGradient>
    //         <linearGradient
    //           id="paint5_linear_139_452"
    //           x1="19"
    //           y1="3"
    //           x2="19"
    //           y2="25"
    //           gradientUnits="userSpaceOnUse"
    //         >
    //           <stop stop-color="#254EDB" />
    //           <stop offset="1" stop-color="#142A75" />
    //         </linearGradient>
    //       </defs>
    //     </svg>
    //   </Box>
    // );

    if (disabledLink) {
      return logo;
    }

    return (
      <Link
        component={RouterLink}
        href={paths.dashboard.station.root}
        sx={{ display: "contents" }}
      >
        {logo}
      </Link>
    );
  }
);

export default Logo;
