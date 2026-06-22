import { Injectable, NotFoundException, Logger } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Example } from './entities/example.entity';
import { CreateExampleDto } from './dtos/create-example.dto';
import { UpdateExampleDto } from './dtos/update-example.dto';

/**
 * Example Service
 * Business logic layer
 * FROM: Ceng-Carpool service pattern
 */
@Injectable()
export class ExampleService {
  private readonly logger = new Logger(ExampleService.name);

  constructor(
    @InjectRepository(Example)
    private readonly exampleRepository: Repository<Example>,
  ) {}

  async findAll(): Promise<Example[]> {
    this.logger.debug('Finding all examples');
    return this.exampleRepository.find();
  }

  async findOne(id: string): Promise<Example> {
    this.logger.debug(`Finding example with ID: ${id}`);
    const example = await this.exampleRepository.findOneBy({ id });

    if (!example) {
      this.logger.warn(`Example not found: ${id}`);
      throw new NotFoundException(`Example with ID ${id} not found`);
    }

    return example;
  }

  async create(createDto: CreateExampleDto): Promise<Example> {
    this.logger.debug(`Creating example: ${JSON.stringify(createDto)}`);
    const example = this.exampleRepository.create(createDto);
    const saved = await this.exampleRepository.save(example);
    this.logger.log(`Example created with ID: ${saved.id}`);
    return saved;
  }

  async update(id: string, updateDto: UpdateExampleDto): Promise<Example> {
    this.logger.debug(`Updating example: ${id}`);
    const example = await this.findOne(id);
    
    Object.assign(example, updateDto);
    const updated = await this.exampleRepository.save(example);
    
    this.logger.log(`Example updated: ${id}`);
    return updated;
  }

  async remove(id: string): Promise<void> {
    this.logger.debug(`Deleting example: ${id}`);
    const example = await this.findOne(id);
    
    await this.exampleRepository.remove(example);
    this.logger.log(`Example deleted: ${id}`);
  }
}
