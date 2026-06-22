import { Controller, Get, Post, Body, Param, Delete, Patch } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';
import { ExampleService } from './example.service';
import { CreateExampleDto } from './dtos/create-example.dto';
import { UpdateExampleDto } from './dtos/update-example.dto';

/**
 * Example Controller
 * Demonstrates CRUD endpoints
 * FROM: Ceng-Carpool controller pattern
 */
@ApiTags('example')
@Controller('example')
export class ExampleController {
  constructor(private readonly exampleService: ExampleService) {}

  @Get()
  @ApiOperation({ summary: 'Get all examples' })
  @ApiResponse({ status: 200, description: 'List of examples' })
  async findAll() {
    return this.exampleService.findAll();
  }

  @Get(':id')
  @ApiOperation({ summary: 'Get example by ID' })
  @ApiResponse({ status: 200, description: 'Example found' })
  @ApiResponse({ status: 404, description: 'Example not found' })
  async findOne(@Param('id') id: string) {
    return this.exampleService.findOne(id);
  }

  @Post()
  @ApiOperation({ summary: 'Create new example' })
  @ApiResponse({ status: 201, description: 'Example created' })
  async create(@Body() createDto: CreateExampleDto) {
    return this.exampleService.create(createDto);
  }

  @Patch(':id')
  @ApiOperation({ summary: 'Update example' })
  @ApiResponse({ status: 200, description: 'Example updated' })
  async update(
    @Param('id') id: string,
    @Body() updateDto: UpdateExampleDto,
  ) {
    return this.exampleService.update(id, updateDto);
  }

  @Delete(':id')
  @ApiOperation({ summary: 'Delete example' })
  @ApiResponse({ status: 200, description: 'Example deleted' })
  async remove(@Param('id') id: string) {
    return this.exampleService.remove(id);
  }
}
