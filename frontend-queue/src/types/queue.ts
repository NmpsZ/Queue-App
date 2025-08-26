export interface Queue {
  id: number;
  queue_no: string;  // ต้องตรงกับ json tag จาก Go
  name: string;
  phone: string;
  status: string;
  qr_code?: string;
  created_at?: string;
  updated_at?: string;
}