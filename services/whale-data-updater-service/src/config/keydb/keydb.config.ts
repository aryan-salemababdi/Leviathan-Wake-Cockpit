import { registerAs } from '@nestjs/config';

export default registerAs('keydb', () => ({
  host: process.env.KEYDB_HOST || 'localhost',
  port: parseInt(process.env.KEYDB_PORT!, 10) || 6379,
}));
