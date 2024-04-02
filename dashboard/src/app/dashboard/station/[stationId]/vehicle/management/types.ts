// ----------------------------------------------------------------------

import { VehicleTargetConfiguration, ImsiConfiguration } from "@/apollo-client";

export type IVehicleTableFilterValue = string | Date | null;

export type IVehicleTableFilters = {
  name: string;
};

// ----------------------------------------------------------------------

export type IVehicleHistory = {
  VehicleTime: Date;
  paymentTime: Date;
  deliveryTime: Date;
  completionTime: Date;
  timeline: {
    title: string;
    time: Date;
  }[];
};

export type IVehicleShippingAddress = {
  fullAddress: string;
  phoneNumber: string;
};

export type IVehiclePayment = {
  cardType: string;
  cardNumber: string;
};

export type IVehicleDelivery = {
  shipBy: string;
  speedy: string;
  trackingNumber: string;
};

export type IVehicleCustomer = {
  id: string;
  name: string;
  email: string;
  avatarUrl: string;
  ipAddress: string;
};

export type IVehicleItem = VehicleTargetConfiguration;
export type IImsiConfigurationItem = ImsiConfiguration;
