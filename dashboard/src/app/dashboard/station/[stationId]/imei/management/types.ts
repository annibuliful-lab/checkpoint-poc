// ----------------------------------------------------------------------

import {
  ImeiConfiguration,
  ImsiConfiguration,
  MobileDeviceConfiguration,
} from "@/apollo-client";

export type IImeiTableFilterValue = string | Date | null;

export type IImeiTableFilters = {
  name: string;
};

// ----------------------------------------------------------------------

export type IImeiHistory = {
  ImeiTime: Date;
  paymentTime: Date;
  deliveryTime: Date;
  completionTime: Date;
  timeline: {
    title: string;
    time: Date;
  }[];
};

export type IImeiShippingAddress = {
  fullAddress: string;
  phoneNumber: string;
};

export type IImeiPayment = {
  cardType: string;
  cardNumber: string;
};

export type IImeiDelivery = {
  shipBy: string;
  speedy: string;
  trackingNumber: string;
};

export type IImeiCustomer = {
  id: string;
  name: string;
  email: string;
  avatarUrl: string;
  ipAddress: string;
};

export type IImeiItem = MobileDeviceConfiguration;
export type IImsiConfigurationItem = ImsiConfiguration;
