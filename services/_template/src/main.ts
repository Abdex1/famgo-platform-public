import { NestFactory } from '@nestjs/common';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { Logger, ValidationPipe } from '@nestjs/common';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  const logger = new Logger('FamGo');

  // Enable CORS
  app.enableCors({
    origin: process.env.CORS_ORIGIN || '*',
    credentials: true,
  });

  // Global validation pipe
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      forbidNonWhitelisted: true,
      transform: true,
      transformOptions: {
        enableImplicitConversion: true,
      },
    }),
  );

  // Swagger API Documentation
  const config = new DocumentBuilder()
    .setTitle('FamGo ' + (process.env.SERVICE_NAME || 'Service'))
    .setDescription(
      'Enterprise Urban Mobility Operating System - ' +
        (process.env.SERVICE_NAME || 'Template Service'),
    )
    .setVersion('1.0')
    .addBearerAuth(
      {
        type: 'http',
        scheme: 'bearer',
        bearerFormat: 'JWT',
        description: 'Enter JWT token',
      },
      'JWT',
    )
    .addServer(process.env.GATEWAY_URL || 'http://localhost:8000', 'Production')
    .addServer('http://localhost:3000', 'Local')
    .build();

  const document = SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('api/docs', app, document);

  // Health check endpoint
  app.get('/health', () => ({
    status: 'ok',
    timestamp: new Date().toISOString(),
    service: process.env.SERVICE_NAME,
  }));

  const port = process.env.PORT || 3000;
  await app.listen(port);

  logger.log(`
╔════════════════════════════════════════════════════════════╗
║     FamGo Platform - Enterprise Mobility Operating System  ║
╠════════════════════════════════════════════════════════════╣
║  Service: ${(process.env.SERVICE_NAME || 'Template Service').padEnd(46)} ║
║  Environment: ${(process.env.NODE_ENV || 'development').padEnd(39)} ║
║  Port: ${String(port).padEnd(54)} ║
║  API Docs: http://localhost:${String(port).padEnd(38)}/api/docs ║
║  Health Check: http://localhost:${String(port).padEnd(33)}/health ║
╚════════════════════════════════════════════════════════════╝
  `);
}

bootstrap().catch((err) => {
  console.error('Failed to start service:', err);
  process.exit(1);
});
