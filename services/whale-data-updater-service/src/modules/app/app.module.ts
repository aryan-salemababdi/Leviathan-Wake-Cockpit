import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import keydbConfig from 'src/config/keydb/keydb.config';
import { ScheduleModule } from '@nestjs/schedule';
import { KeydbModule } from '../keydb/keydb.module';
import { UpdaterModule } from '../updater/updater.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: [keydbConfig],
    }),
    ScheduleModule.forRoot(),
    KeydbModule,
    UpdaterModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}
