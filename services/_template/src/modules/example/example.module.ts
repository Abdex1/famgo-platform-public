import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ExampleController } from './example.controller';
import { ExampleService } from './example.service';
import { Example } from './entities/example.entity';

/**
 * Example Module
 * Template for creating new feature modules
 * FROM: Ceng-Carpool module pattern
 */
@Module({
  imports: [TypeOrmModule.forFeature([Example])],
  controllers: [ExampleController],
  providers: [ExampleService],
  exports: [ExampleService], // Export if other modules need this service
})
export class ExampleModule {}
