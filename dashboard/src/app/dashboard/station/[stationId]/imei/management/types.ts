// ----------------------------------------------------------------------

import { ImsiConfiguration, MobileDeviceConfiguration } from "@/apollo-client";

export type ImeiTableFilterValue = string | Date | null;

export type ImeiTableFilters = {
  name: string;
};

// ----------------------------------------------------------------------

export type ImeiHistory = {
  ImeiTime: Date;
  paymentTime: Date;
  deliveryTime: Date;
  completionTime: Date;
  timeline: {
    title: string;
    time: Date;
  }[];
};

export type ImeiShippingAddress = {
  fullAddress: string;
  phoneNumber: string;
};

export type ImeiPayment = {
  cardType: string;
  cardNumber: string;
};

export type ImeiDelivery = {
  shipBy: string;
  speedy: string;
  trackingNumber: string;
};

export type ImeiCustomer = {
  id: string;
  name: string;
  email: string;
  avatarUrl: string;
  ipAddress: string;
};

export type ImeiItem = MobileDeviceConfiguration;
export type IImsiConfigurationItem = ImsiConfiguration;
