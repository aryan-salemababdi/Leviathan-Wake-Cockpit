import { Module } from "@nestjs/common";
import { KeydbModule } from "../keydb/keydb.module";
import { HttpModule } from "@nestjs/axios";
import { UpdaterService } from "./updater.service";


@Module({
    imports: [KeydbModule, HttpModule],
    providers: [UpdaterService]
})

export class UpdaterModule {};