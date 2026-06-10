export interface AuthStatus {
  connected: boolean;
  merchantId?: string;
  name?: string;
  cloverEnv?: string;
}

export interface BootstrapRequest {
  name: string;
  clover_merchant_id: string;
  access_token: string;
  refresh_token?: string;
}

export interface BootstrapResponse {
  id: string;
  name: string;
  clover_merchant_id: string;
  is_connected: boolean;
}

export interface RevokeRequest {
  merchant_id: string;
}
