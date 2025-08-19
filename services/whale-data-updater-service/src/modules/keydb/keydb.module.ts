import { Module } from '@nestjs/common';
import { KeydbService } from './keydb.service';

@Module({
  providers: [KeydbService],
  exports: [KeydbService],
})
export class KeydbModule {}
